package main

import (
	"encoding/csv"
	"fmt"

	log "github.com/Sirupsen/logrus"
	// considering using this one instead : i18n4go, tools better but not inside
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/nicksnyder/go-i18n/i18n"
	stockCmd "github.com/olivier5741/stock-manager/cmd/stock"
	"github.com/olivier5741/stock-manager/item"
	"github.com/olivier5741/stock-manager/item/amount"
	"github.com/olivier5741/stock-manager/item/items"
	"github.com/olivier5741/stock-manager/item/quant"
	"github.com/olivier5741/stock-manager/skelet"
	"github.com/olivier5741/stock-manager/stock"
	"github.com/olivier5741/strtab"
)

// ATTENTION : tout méthode qui modifie le struct doit accepter un pointer !! ??
// Windows build :
// GOOS=windows GOARCH=386 go build -o stock-manager-0.1.exe main.go config.go
var (
	Tr i18n.TranslateFunc

	sep = "-"

	csvSuff    = ".csv"
	TimeFormat = "2006" + sep + "01" + sep + "02"

	extension       = ".csv"
	configPrefix    = "c" + sep
	exePrefix       = "e" + sep
	generatedPrefix = "g" + sep
	loggingPrefix   = "l" + sep
	numberPrefix    = "n"

	repo  = stockCmd.MakeDummyStockRepo()
	endPt = stockCmd.EndPt{Db: repo}

	stockRoute = func(t skelet.Ider) (ok bool, a skelet.AggAct, p skelet.EvtSrcPersister) {
		switch t.(type) {
		case stock.InCmd:
			return true, endPt.HandleIn, repo
		case stock.OutCmd:
			return true, endPt.HandleOut, repo
		case stock.InventoryCmd:
			return true, endPt.HandleInventory, repo
		default:
			return false, nil, nil
		}
	}
)

func csvToStruct(filename string, h []string, mapper func(s []string, c interface{}),
	newLiner func() interface{}, appender func(interface{})) {

	file, err2 := os.OpenFile(filename, os.O_CREATE, 0666)
	defer file.Close()
	if err2 != nil {
		log.WithFields(log.Fields{
			"filename": filename,
			"err":      err2,
		}).Error(Tr("error_file_open"))
	}

	r := csv.NewReader(file)
	out, err2a := r.ReadAll()

	if err2a != nil {
		log.WithFields(log.Fields{
			"filename": filename,
			"err":      err2a,
		}).Error(Tr("error_file_csv_unmarshal"))
	}

	if len(out) < 2 {
		return
	}

	for _, line := range out[1:] {
		newLine := newLiner()
		mapper(line, newLine)
		appender(newLine)
	}
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
	f, err := os.Create(t.Title)

	if err != nil {
		log.WithFields(log.Fields{
			"path":     t.Path,
			"filename": t.Title,
			"err":      err,
		}).Error("create_file_error")
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
		}).Error(Tr("error_file_csv_view_to"))
	}
}

func init() {
	i18n.MustLoadTranslationFile("c-int/en-us.all.yaml")
	i18n.MustLoadTranslationFile("c-int/fr-be.all.yaml")
	Tr, _ = i18n.Tfunc("fr-be")
}

func main() {

	// LOGGING
	logfile := loggingPrefix + Tr("file_name_log")
	f, err1 := os.OpenFile(logfile,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err1 != nil {
		log.WithFields(log.Fields{
			"filename": logfile,
			"err":      err1,
		}).Error(Tr("error_file_open"))
	}
	defer f.Close()

	log.SetOutput(f)
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.ErrorLevel)

	files, err3 := ioutil.ReadDir("./")

	if err3 != nil {
		log.WithFields(log.Fields{
			"err": err3,
		}).Error(Tr("error_dir_read"))
	}

	// DRIVE DOWNLOAD
	//getFiles()

	RouteFile(files)

	stockInt, err4 := endPt.Db.Get("main")
	iStock := stockInt.(*stock.Stock).Items
	if err4 != nil {
		log.WithFields(log.Fields{
			"err": err4,
		}).Error(Tr("error_query_stock")) // Cannot execute stock query
	}

	prodValHeader := []string{Tr("csv_header_item_product"),
		Tr("csv_header_item_value", 1), Tr("csv_header_item_unit", 1),
		Tr("csv_header_item_value", 2), Tr("csv_header_item_unit", 2),
		Tr("csv_header_item_value", 3), Tr("csv_header_item_unit", 3),
		Tr("csv_header_item_value", 4), Tr("csv_header_item_unit", 4),
	}
	prodValRender := func(tab *strtab.T) [][]string {
		return tab.GetContentWithHeaders(false)
	}
	prodEvolRender := func(tab *strtab.T) [][]string {
		return tab.GetContentWithHeaders(true)
	}

	TableView{"main", generatedPrefix + Tr("file_name_stock") + extension,
		strtab.NewT(prodValHeader, iStock.StringSlice()...).Sort(),
		prodValRender}.Show()

	TableView{"main", draftFileName(3, Tr("file_name_inventory")),
		strtab.NewT(prodValHeader, iStock.StringSlice()...).Sort(),
		prodValRender}.Show()

	TableView{"main", draftFileName(1, Tr("file_name_stock_in")),
		strtab.NewT(prodValHeader, iStock.Empty().StringSlice()...).Sort(),
		prodValRender}.Show()

	TableView{"main", draftFileName(2, Tr("file_name_stock_out")),
		strtab.NewT(prodValHeader, iStock.Empty().StringSlice()...).Sort(),
		prodValRender}.Show()

	TableView{"main", generatedPrefix + Tr("file_name_product") + extension,
		strtab.NewTfromMap(items.ItemsMapToStringMapTable(
			endPt.ProdValEvol("main"))).Sort().Transpose().Sort(),
		prodEvolRender}.Show()
}

func draftFileName(num int, name string) string {
	today := time.Now().Format(TimeFormat)
	return Tr("file_prefix_draft") + sep + today + sep + numberPrefix +
		strconv.Itoa(num) + sep + name + extension
}

type FileInputRouter struct {
	Routes map[string]func(agg interface{}, cmd interface{}) (event interface{}, extEvent interface{}, err error)
}

func RouteFile(files []os.FileInfo) {
	for _, file := range files {
		// TO REFACTOR

		if !strings.HasSuffix(file.Name(), ".csv") {
			continue
		}

		if strings.HasPrefix(file.Name(), generatedPrefix) ||
			strings.HasPrefix(file.Name(), Tr("file_prefix_draft")) {
			continue
		}

		log.Debug(file.Name())

		path, err1 := ParseFilename(file.Name())
		if err1 != nil {
			log.WithFields(log.Fields{
				"filename": file.Name(),
				"err":      err1,
			}).Error(Tr("error_filename_parse")) //Cannot parse filename
			continue
		}

		out, err3 := UnmarshalCsvFile(path)
		if err3 != nil {
			log.WithFields(log.Fields{
				"filename": file.Name(),
				"err":      err3,
			}).Error(Tr("error_file_csv_unmarshal")) //"Cannot unmarshal to csv"
			continue
		}

		skelet.ExecuteCommand(skelet.Cmd{T: out, Route: stockRoute}, stockCmd.Chain)
	}
}

func UnmarshalCsvFile(path Filename) (out skelet.Ider, err error) {

	// should be somewhere else perhaps
	mapper := func(ins []string, c interface{}) {
		c.(*item.Item).Prod = item.Prod(ins[0])
		var units []quant.Quant
		for i := 1; i < len(ins)-1; i = i + 2 {
			val, _ := strconv.Atoi(ins[i])
			log.Debug("Val")
			log.Debug(val)
			units = append(units, quant.Quant{quant.NewUnit(ins[i+1]), val})
			log.Debug(units)
		}
		c.(*item.Item).Amount = amount.NewAmount(units...)
	}

	// should put this in a local type
	headers := []string{
		Tr("csv_header_item_product"),
		Tr("csv_header_item_value", 1),
		Tr("csv_header_item_unit", 1),
		Tr("csv_header_item_value", 2),
		Tr("csv_header_item_unit", 2),
		Tr("csv_header_item_value", 3),
		Tr("csv_header_item_unit", 3),
		Tr("csv_header_item_value", 4),
		Tr("csv_header_item_unit", 4),
	}

	var its []item.Item
	newLiner := func() interface{} { return new(item.Item) }
	appender := func(v interface{}) {
		a := v.(*item.Item)
		its = append(its, *a)
	}

	// should put this inside switch case
	csvToStruct(path.String(), headers, mapper, newLiner, appender)
	itsMap := items.FromSlice(its)

	switch path.Act {
	case Tr("file_name_stock_in"):
		return stock.InCmd{path.Stock, itsMap, path.Date.Format(TimeFormat) + sep + numberPrefix + path.ID}, nil
	case Tr("file_name_stock_out"):
		return stock.OutCmd{path.Stock, itsMap, path.Date.Format(TimeFormat) + sep + numberPrefix + path.ID}, nil
	case Tr("file_name_inventory"):
		return stock.InventoryCmd{path.Stock, itsMap, path.Date.Format(TimeFormat) + sep + numberPrefix + path.ID}, nil

	}
	return nil, fmt.Errorf(Tr("no_action_for_filename_error")) // No action found
}

func ParseFilename(s string) (f Filename, err error) {
	f = Filename{}

	s = filepath.Base(s)
	//Extension
	f.Ext = filepath.Ext(s)
	if strings.HasSuffix(s, csvSuff) {
		err = fmt.Errorf(Tr("filename_csv_suffix_error"), s)
	}

	ts := strings.TrimSuffix(s, f.Ext)

	args := strings.Split(ts, sep)
	if l := len(args); l != 5 && l != 6 {
		err = fmt.Errorf(Tr("filename_number_argument_error"), s)
		return
	}

	if len(args) == 6 {
		f.Status = args[0]
		args = args[0:]
	}

	log.Debug(args)

	//Date
	date, err := time.Parse(TimeFormat, args[0]+sep+args[1]+sep+args[2])
	if err != nil {
		return
	}
	f.Date = date

	//Stock
	//f.Stock = args[3]
	f.Stock = "main" //To refactor

	//Id
	f.ID = strings.TrimPrefix(args[3], numberPrefix)

	//Action
	f.Act = args[4]

	return
}

type Filename struct {
	Date                        time.Time
	Stock, Act, ID, Ext, Status string
}

func (f Filename) String() string {
	s := f.Date.Format(TimeFormat) + sep + numberPrefix + f.ID + sep + f.Act
	if f.Status != "" {
		s = s + sep + f.Status
	}
	s = s + f.Ext

	return s
}
