package main

import (
	"encoding/csv"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gocarina/gocsv"
	"github.com/olivier5741/stock-manager/cmd/stock"
	. "github.com/olivier5741/stock-manager/item"
	"github.com/olivier5741/stock-manager/skelet"
	. "github.com/olivier5741/stock-manager/stock"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

// ATTENTION : tout méthode qui modifie le struct doit accepter un pointer !! ??
// Windows build :
// GOOS=windows GOARCH=386 go build -o stock-manager-0.1.exe main.go config.go
var (
	dir = "Bièvre"

	extension       = ".csv"
	configPrefix    = "c-"
	generatedPrefix = "g-"
	loggingPrefix   = "l-"
	draftSuffix     = "-en attente" + extension
	numberPrefix    = "n°"

	configFilename = configPrefix + "config.csv"
	logfile        = loggingPrefix + "erreurs"
	dashboardfile  = generatedPrefix + "stock" + extension
	orderfile      = generatedPrefix + "à commander" + extension

	inventoryFilename = "inventaire"
	inFilename        = "entrée"
	outFilename       = "sortie"
	orderFilename     = "commande"

	cannotOpenFile = "Cannot open file"

	repo  = stock.MakeDummyStockRepository()
	endPt = stock.EndPt{Db: repo}

	stockRoute = func(t skelet.Ider) (ok bool, a skelet.AggAct, p skelet.EvtSrcPersister) {
		switch t.(type) {
		case skelet.InCmd:
			return true, endPt.HandleIn, repo
		case skelet.OutCmd:
			return true, endPt.HandleOut, repo
		case skelet.InventoryCmd:
			return true, endPt.HandleInventory, repo
		default:
			return false, nil, nil
		}
	}
)

func itemArrayToMap(items []Item) (out Items) {
	out = Items{}
	for _, item := range items {
		out[string(item.Prod)] = item
	}
	return
}

func logIfCannotOpenFile(err error, f string) {
	if err != nil {
		log.WithFields(log.Fields{
			"filename": logfile,
		}).Error(cannotOpenFile)
	}
}

func GetConfigFromFile(filename string, config interface{}) error {
	f, err1 := os.OpenFile(filename, os.O_CREATE, 0666)
	if err1 != nil {
		return err1
	}
	defer f.Close()

	out, err2 := ioutil.ReadAll(f)
	if err2 != nil {
		return err2
	}
	return yaml.Unmarshal(out, config)
}

func main() {

	// setup environment if not exist
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, 0755)
	}

	f, err1 := os.OpenFile(dir+"/"+logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err1 != nil {
		log.WithFields(log.Fields{
			"filename": logfile,
			"err":      err1,
		}).Error(cannotOpenFile)
	}
	defer f.Close()

	log.SetOutput(f)
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.DebugLevel)

	config := make([]ConfigProd, 0)
	configFile, err2 := os.OpenFile(dir+"/"+configFilename, os.O_CREATE, 0666)
	defer configFile.Close()
	if err2 != nil {
		log.WithFields(log.Fields{
			"filename": configFilename,
			"err":      err2,
		}).Error(cannotOpenFile)
	}

	err2a := gocsv.UnmarshalFile(configFile, &config)

	if err2a != nil {
		log.WithFields(log.Fields{
			"filename": configFilename,
			"err":      err2a,
		}).Error("Cannot get config from csv file")
	}

	files, err3 := ioutil.ReadDir(dir)

	if err3 != nil {
		log.WithFields(log.Fields{
			"dir": dir,
			"err": err3,
		}).Error("Cannot read directory")
	}

	RouteFile(files)

	// OUTPUT
	// Left in stock
	p := stock.MakeProdInStockTable()
	r, err4 := endPt.StocksQuery()
	inventory := r[0].Items.Copy()
	if err4 != nil {
		log.WithFields(log.Fields{
			"err": err4,
		}).Error("Cannot execute stock query")
	}

	p.Parse(r)
	err5 := WriteCsvFile(p.ToProductStringLines(), dashboardfile)
	if err5 != nil {
		log.WithFields(log.Fields{
			"filename": dashboardfile,
			"err":      err5,
		}).Error("Cannot write csv file")
	}

	//Missing
	mt := stock.MakeProdInStockTable()
	m := make([]*Stock, 0)
	for _, s := range r {
		s.Items = s.Items.Missing(GetMissingItems(config))
		m = append(m, s)
	}
	order := m[0].Items.Copy()
	mt.Parse(m)
	err6 := WriteCsvFile(mt.ToProductStringLines(), orderfile)
	if err6 != nil {
		log.WithFields(log.Fields{
			"filename": orderfile,
			"err":      err6,
		}).Error("Cannot write csv file")
	}

	// Create missing files
	today := time.Now().Format("2006-01-02")
	inDraftFilename := today + "-" + numberPrefix + "1-" + inFilename + draftSuffix
	outDraftFilename := today + "-" + numberPrefix + "2-" + outFilename + draftSuffix
	invDraftFilename := today + "-" + numberPrefix + "3-" + inventoryFilename + draftSuffix
	orderDraftFilename := today + "-" + numberPrefix + "4-" + orderFilename + draftSuffix

	log.Debug(m[0].Items.ToStringLines())

	createDateOrUpdateDate(inDraftFilename, today, addHeader(ToItemStringLines(config)))
	createDateOrUpdateDate(outDraftFilename, today, addHeader(ToItemStringLines(config)))
	createDateOrUpdateDate(invDraftFilename, today, addHeader(inventory.ToStringLines()))
	createDateOrUpdateDate(orderDraftFilename, today, addHeader(order.ToStringLines()))
}

func addHeader(ins [][]string) [][]string {
	out := [][]string{[]string{"Prod", "Val"}}
	for _, in := range ins {
		out = append(out, in)
	}
	return out
}

func createDateOrUpdateDate(filename string, today string, lines [][]string) {
	filename = dir + "/" + filename
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		f, err1 := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err1 != nil {
			log.WithFields(log.Fields{
				"filename": filename,
				"err":      err1,
			}).Error(cannotOpenFile)
		}

		defer f.Close()

		w := csv.NewWriter(f)
		w.WriteAll(lines)
		if w.Error() != nil {
			log.WithFields(log.Fields{
				"filename": filename,
				"err":      w.Error(),
			}).Error("Cannot write csv to file")
		}
	}
}

type FileInputRouter struct {
	Routes map[string]func(agg interface{}, cmd interface{}) (event interface{}, extEvent interface{}, err error)
}

func RouteFile(files []os.FileInfo) {
	for _, file := range files {
		// TO REFACTOR
		if strings.HasPrefix(file.Name(), configPrefix) ||
			strings.HasPrefix(file.Name(), generatedPrefix) ||
			strings.HasPrefix(file.Name(), loggingPrefix) ||
			strings.HasSuffix(file.Name(), draftSuffix) {
			return
		}
		path, err1 := ParseFilename(file.Name())
		if err1 != nil {
			log.WithFields(log.Fields{
				"filename": file.Name(),
				"err":      err1,
			}).Error("Cannot parse filename")
			continue
		}

		out, err3 := UnmarshalCsvFile(path)
		if err3 != nil {
			log.WithFields(log.Fields{
				"filename": file.Name(),
				"err":      err3,
			}).Error("Cannot unmarshal to csv")
			continue
		}

		skelet.ExecuteCommand(skelet.Cmd{T: out, Route: stockRoute}, stock.Chain)
	}
}

func WriteCsvFile(lines [][]string, path string) error {
	f, err := os.Create(dir + "/" + path)
	defer f.Close()
	if err != nil {
		return err
	}
	w := csv.NewWriter(f)
	//w.Comma = ';'

	w.WriteAll(lines)
	return w.Error()
}

func UnmarshalCsvFile(path Filename) (out skelet.Ider, err error) {
	f, err := os.Open(dir + "/" + path.String())
	defer f.Close()
	if err != nil {
		return nil, err
	}

	gocsv.SetCSVReader(func(out io.Reader) *csv.Reader {
		r := csv.NewReader(out)
		//r.Comma = ';'
		return r
	})

	switch path.Act {
	case inFilename:
		its := []Item{}
		err := gocsv.UnmarshalFile(f, &its)
		if err != nil {
			return nil, err
		}
		return skelet.InCmd{path.Stock, itemArrayToMap(its)}, nil
	case outFilename:
		its := []Item{}
		err := gocsv.UnmarshalFile(f, &its)
		if err != nil {
			return nil, err
		}
		return skelet.OutCmd{path.Stock, itemArrayToMap(its)}, nil
	case inventoryFilename:
		its := []Item{}
		err := gocsv.UnmarshalFile(f, &its)
		if err != nil {
			return nil, err
		}
		return skelet.InventoryCmd{path.Stock, itemArrayToMap(its)}, nil

	}
	return nil, fmt.Errorf("No action found")
}
