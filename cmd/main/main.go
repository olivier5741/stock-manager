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
)

// ATTENTION : tout méthode qui modifie le struct doit accepter un pointer !! ??
// Windows build :
// GOOS=windows GOARCH=386 go build -o stock-manager-0.1.exe main.go config.go
var (
	inputDir       = "input"
	outputDir      = "output"
	yamlConfigFile = "config.yaml"
	logfile        = "log"
	dashboardfile  = "stock.csv"
	orderfile      = "missing.csv"

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
	if _, err := os.Stat(inputDir); os.IsNotExist(err) {
		os.Mkdir(inputDir, 0755)
	}

	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		os.Mkdir(outputDir, 0755)
	}

	f, err1 := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
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

	config := Config{}
	err2 := GetConfigFromFile(yamlConfigFile, &config)
	if err2 != nil {
		log.WithFields(log.Fields{
			"filename": yamlConfigFile,
			"err":      err2,
		}).Error("Cannot get config from yaml file")
	}

	//log.Debug(config)

	files, err3 := ioutil.ReadDir(inputDir)
	if err3 != nil {
		log.WithFields(log.Fields{
			"dir": inputDir,
			"err": err3,
		}).Error("Cannot read directory")
	}

	RouteFile(files)

	// OUTPUT
	// Left in stock
	p := stock.MakeProdInStockTable()
	r, err4 := endPt.StocksQuery()
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
		s.Items = s.Items.Missing(config.GetMissingItems())
		m = append(m, s)
	}
	mt.Parse(m)
	err6 := WriteCsvFile(mt.ToProductStringLines(), orderfile)
	if err6 != nil {
		log.WithFields(log.Fields{
			"filename": orderfile,
			"err":      err6,
		}).Error("Cannot write csv file")
	}

}

type FileInputRouter struct {
	Routes map[string]func(agg interface{}, cmd interface{}) (event interface{}, extEvent interface{}, err error)
}

func RouteFile(files []os.FileInfo) {
	for _, file := range files {
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
	f, err := os.Create(outputDir + "/" + path)
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
	f, err := os.Open(inputDir + "/" + path.String())
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
	case "in":
		its := []Item{}
		err := gocsv.UnmarshalFile(f, &its)
		if err != nil {
			return nil, err
		}
		return skelet.InCmd{path.Stock, itemArrayToMap(its)}, nil
	case "out":
		its := []Item{}
		err := gocsv.UnmarshalFile(f, &its)
		if err != nil {
			return nil, err
		}
		return skelet.OutCmd{path.Stock, itemArrayToMap(its)}, nil
	case "inv":
		its := []Item{}
		err := gocsv.UnmarshalFile(f, &its)
		if err != nil {
			return nil, err
		}
		return skelet.InventoryCmd{path.Stock, itemArrayToMap(its)}, nil

	}
	return nil, fmt.Errorf("No action found")
}
