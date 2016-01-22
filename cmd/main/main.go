package main

import (
	"encoding/csv"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gocarina/gocsv"
	"github.com/olivier5741/stock-manager/cmd/stock"
	. "github.com/olivier5741/stock-manager/item"
	"github.com/olivier5741/stock-manager/skelet"
	stockBL "github.com/olivier5741/stock-manager/stock"
	"github.com/olivier5741/stock-manager/strtab"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

// ATTENTION : tout méthode qui modifie le struct doit accepter un pointer !! ??
// Windows build :
// GOOS=windows GOARCH=386 go build -o stock-manager-0.1.exe main.go config.go
var (
	dir = "bievre"

	extension       = ".csv"
	configPrefix    = "c-"
	generatedPrefix = "g-"
	loggingPrefix   = "l-"
	draftPrefix     = "brouillon"
	numberPrefix    = "n°"

	configFilename = configPrefix + "config.csv"
	logfile        = loggingPrefix + "erreurs"

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

func getConfig() []ConfigProd {
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

	return config
}

func inStock(config []ConfigProd) (its Items) {
	stock1, err4 := endPt.Db.Get("bievre")
	if err4 != nil {
		log.WithFields(log.Fields{
			"err": err4,
		}).Error("Cannot execute stock query")
	}

	its = stock1.(*stockBL.Stock).Items.Copy()

	for _, prod := range config {
		if _, ok := its[prod.Name]; !ok {
			its[prod.Name] = Item{Prod(prod.Name), Val{0}}
		}
	}

	return
}

func missing(config []ConfigProd, s Items) Items {
	return s.Missing(GetMissingItems(config))
}

func mapItem(its Items) [][]string {
	out := make([][]string, 0)
	for _, it := range its {
		out = append(out, []string{it.Prod.String(), it.Val.String()})
	}
	return out
}

func mapConfigProd(cs []ConfigProd) [][]string {
	out := make([][]string, 0)
	for _, c := range cs {
		out = append(out, []string{c.Name, ""})
	}
	return out
}

func mapItemsMap(its map[string]Items) map[string]map[string]string {
	out := make(map[string]map[string]string, 0)
	for date, it := range its {
		newRow := make(map[string]string)
		for prod, val := range it {
			newRow[prod] = val.Val.String()
		}
		out[date] = newRow
	}
	return out
}

type Viewer interface {
	Show()
}

type TableView struct {
	Path   string
	Title  string
	Table  *strtab.T
	Render func(*strtab.T) [][]string
}

func (t TableView) Show() {
	f, err := os.Create(t.Path + "/" + t.Title)

	if err != nil {
		log.WithFields(log.Fields{
			"path":     t.Path,
			"filename": t.Title,
			"err":      err,
		}).Error("Error creating the csv file")
		return
	}

	defer f.Close()

	w := csv.NewWriter(f)
	w.WriteAll(t.Render(t.Table))

	if w.Error() != nil {
		log.WithFields(log.Fields{
			"path":     t.Path,
			"filename": t.Title,
			"view":     t.Table,
			"err":      w.Error(),
		}).Error("Error writing the view to a csv file")
	}
}

func main() {

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
	log.SetLevel(log.ErrorLevel)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, 0755)
	}

	config := getConfig()

	files, err3 := ioutil.ReadDir(dir)

	if err3 != nil {
		log.WithFields(log.Fields{
			"dir": dir,
			"err": err3,
		}).Error("Cannot read directory")
	}

	RouteFile(files)

	iStock := inStock(config).Copy()

	prodValHeader := []string{"Prod", "Val"}
	prodValRender := func(tab *strtab.T) [][]string { return tab.GetContentWithHeaders(false) }
	prodEvolRender := func(tab *strtab.T) [][]string { return tab.GetContentWithHeaders(true) }

	TableView{"bievre", "g-stock.csv",
		strtab.NewTable(prodValHeader, mapItem(iStock)...), prodValRender}.Show()

	TableView{"bievre", draftFileName(3, inventoryFilename),
		strtab.NewTable(prodValHeader, mapItem(iStock)...), prodValRender}.Show()

	//Missing
	missing := missing(config, iStock)

	TableView{"bievre", "g-à commander.csv",
		strtab.NewTable(prodValHeader, mapItem(missing)...), prodValRender}.Show()

	TableView{"bievre", draftFileName(4, orderFilename),
		strtab.NewTable(prodValHeader, mapItem(missing)...), prodValRender}.Show()

	TableView{"bievre", draftFileName(1, inFilename),
		strtab.NewTable(prodValHeader, mapConfigProd(config)...), prodValRender}.Show()

	TableView{"bievre", draftFileName(2, outFilename),
		strtab.NewTable(prodValHeader, mapConfigProd(config)...), prodValRender}.Show()

	TableView{"bievre", "g-produits.csv",
		strtab.NewTableFromMap(mapItemsMap(endPt.ProdValEvolution("bievre"))).Transpose(), prodEvolRender}.Show()
}

func draftFileName(num int, name string) string {
	today := time.Now().Format("2006-01-02")
	return draftPrefix + "-" + today + "-" + numberPrefix + strconv.Itoa(num) + "-" + name + extension
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
			strings.HasPrefix(file.Name(), draftPrefix) {
			continue
		}

		log.Debug(file.Name())

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
		return skelet.InCmd{path.Stock, itemArrayToMap(its), path.Date.Format(TimeFormat) + "-" + numberPrefix + path.Id}, nil
	case outFilename:
		its := []Item{}
		err := gocsv.UnmarshalFile(f, &its)
		if err != nil {
			return nil, err
		}
		return skelet.OutCmd{path.Stock, itemArrayToMap(its), path.Date.Format(TimeFormat) + "-" + numberPrefix + path.Id}, nil
	case inventoryFilename:
		its := []Item{}
		err := gocsv.UnmarshalFile(f, &its)
		if err != nil {
			return nil, err
		}
		return skelet.InventoryCmd{path.Stock, itemArrayToMap(its), path.Date.Format(TimeFormat) + "-" + numberPrefix + path.Id}, nil

	}
	return nil, fmt.Errorf("No action found")
}
