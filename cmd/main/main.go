package main

import (
	"encoding/csv"
	"fmt"
	"github.com/gocarina/gocsv"
	"github.com/olivier5741/stock-manager/cmd/stock"
	. "github.com/olivier5741/stock-manager/item"
	. "github.com/olivier5741/stock-manager/skelet"
	"io/ioutil"
	"log"
	"os"
)

// ATTENTION : tout méthode qui modifie le struct doit accepter un pointer !! ??

var (
	inputDir   = "input"
	outputDir  = "output"
	yamlConfig = "config.yaml"

	repo  = stock.MakeDummyStockRepository()
	endPt = stock.EndPt{Db: repo}

	stockRoute = func(t Ider) (ok bool, a AggAct, p EvtSrcPersister) {
		switch t.(type) {
		case InCmd:
			return true, endPt.HandleIn, repo
		case OutCmd:
			return true, endPt.HandleOut, repo
		case InventoryCmd:
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

func logIfErr(err error) {
	if err != nil {
		log.Println(err)
		return
	}
}

func main() {

	config := Config{}
	err := GetConfigFromFile(yamlConfig, &config)
	logIfErr(err)

	files, err1 := ioutil.ReadDir(inputDir)
	logIfErr(err1)

	RouteFile(files)

	// OUTPUT
	// Left in stock
	p := stock.MakeProdInStockTable()
	r, err2 := endPt.StocksQuery()
	logIfErr(err2)

	p.Parse(r)
	err3 := WriteCsvFile(p.ToProductStringLines(), "stock1.csv")
	logIfErr(err3)
	// To order

}

type FileInputRouter struct {
	Routes map[string]func(agg interface{}, cmd interface{}) (event interface{}, extEvent interface{}, err error)
}

func RouteFile(files []os.FileInfo) {
	for _, file := range files {
		path, err1 := ParseFilename(file.Name())
		if err1 != nil {
			log.Println(err1)
			continue
		}

		out, err3 := UnmarshalCsvFile(path)
		if err3 != nil {
			log.Println(err3)
			continue
		}

		ExecuteCommand(Cmd{T: out, Route: stockRoute}, stock.Chain)
	}
}

func WriteCsvFile(lines [][]string, path string) error {
	f, err := os.Create(outputDir + "/" + path)
	defer f.Close()
	if err != nil {
		return err
	}
	w := csv.NewWriter(f)
	w.WriteAll(lines)
	return w.Error()
}

func UnmarshalCsvFile(path Filename) (out Ider, err error) {
	f, err := os.Open(inputDir + "/" + path.String())
	defer f.Close()
	if err != nil {
		return nil, err
	}

	switch path.Act {
	case "in":
		its := []Item{}
		err := gocsv.UnmarshalFile(f, &its)
		if err != nil {
			panic(err)
		}
		return InCmd{path.Stock, itemArrayToMap(its)}, nil
	case "out":
		its := []Item{}
		err := gocsv.UnmarshalFile(f, &its)
		if err != nil {
			panic(err)
		}
		return OutCmd{path.Stock, itemArrayToMap(its)}, nil
	case "inv":
		its := []Item{}
		err := gocsv.UnmarshalFile(f, &its)
		if err != nil {
			panic(err)
		}
		return InventoryCmd{path.Stock, itemArrayToMap(its)}, nil

	}
	return nil, fmt.Errorf("No action found")
}
