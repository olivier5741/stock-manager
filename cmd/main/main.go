package main

import (
	"encoding/csv"
	"fmt"
	log "github.com/Sirupsen/logrus"
	// considering using this one instead : i18n4go, tools better but not inside
	"github.com/nicksnyder/go-i18n/i18n"
	"github.com/olivier5741/stock-manager/cmd/stock"
	. "github.com/olivier5741/stock-manager/item"
	"github.com/olivier5741/stock-manager/skelet"
	stockBL "github.com/olivier5741/stock-manager/stock/main"
	stockSk "github.com/olivier5741/stock-manager/stock/skelet"
	"github.com/olivier5741/stock-manager/strtab"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// ATTENTION : tout méthode qui modifie le struct doit accepter un pointer !! ??
// Windows build :
// GOOS=windows GOARCH=386 go build -o stock-manager-0.1.exe main.go config.go
var (
	Tr i18n.TranslateFunc

	sep = "-"

	csvSuff    = ".csv"
	TimeFormat = "2006" + sep + "01" + sep + "02"

	dir = "bievre"

	extension       = ".csv"
	configPrefix    = "c" + sep
	generatedPrefix = "g" + sep
	loggingPrefix   = "l" + sep
	numberPrefix    = "n°"

	base = Unit{"Base", 1}

	repo  = stock.MakeDummyStockRepository()
	endPt = stock.EndPt{Db: repo}

	stockRoute = func(t skelet.Ider) (ok bool, a skelet.AggAct, p skelet.EvtSrcPersister) {
		switch t.(type) {
		case stockSk.InCmd:
			return true, endPt.HandleIn, repo
		case stockSk.OutCmd:
			return true, endPt.HandleOut, repo
		case stockSk.InventoryCmd:
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

func csvToStruct(filename string, mapper map[string]func(s string, c interface{}),
	newLiner func() interface{}, appender func(interface{})) {

	file, err2 := os.OpenFile(dir+"/"+filename, os.O_CREATE, 0666)
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

	headers := make(map[int]func(string, interface{}), 0)

	for i, h := range out[0] {
		if fun, ok := mapper[h]; ok {
			headers[i] = fun
		}
	}

	for _, line := range out[1:] {
		newLine := newLiner()
		for k, fun := range headers {
			fun(line[k], newLine)
		}
		appender(newLine)
	}
}

func inStock(config []ConfigProd) (its Items) {
	stock1, err4 := endPt.Db.Get("bievre")
	if err4 != nil {
		log.WithFields(log.Fields{
			"err": err4,
		}).Error(Tr("error_query_stock")) // Cannot execute stock query
	}

	its = stock1.(*stockBL.Stock).Items.Copy()

	log.Debug(stock1)

	for _, prod := range config {
		if _, ok := its[prod.Prod]; !ok {
			its[prod.Prod] = Item{Prod(prod.Prod), NewVal(UnitVal{base, 0})}
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
		out = append(out, []string{it.Prod.String(), strconv.Itoa(it.Val.TotalWith())})
	}
	return out
}

func mapConfigProd(cs []ConfigProd) [][]string {
	out := make([][]string, 0)
	for _, c := range cs {
		out = append(out, []string{c.Prod, ""})
	}
	return out
}

func mapItemsMap(its map[string]Items) map[string]map[string]string {
	out := make(map[string]map[string]string, 0)
	for date, it := range its {
		newRow := make(map[string]string)
		for prod, val := range it {
			newRow[prod] = strconv.Itoa(val.Val.TotalWith())
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
		}).Error("create_file_error") //Error creating the csv file
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
		}).Error(Tr("error_file_csv_view_to")) //Error writing the view to a csv file
	}
}

func init() {
	i18n.MustLoadTranslationFile("en-us.all.yaml")
	i18n.MustLoadTranslationFile("fr-be.all.yaml")
	Tr, _ = i18n.Tfunc("fr-be")
}

func main() {

	// LOGGING
	logfile := loggingPrefix + Tr("file_name_log")
	f, err1 := os.OpenFile(dir+"/"+logfile,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err1 != nil {
		log.WithFields(log.Fields{
			"filename": logfile,
			"err":      err1,
		}).Error(Tr("error_file_open"))
	}
	defer f.Close()

	//log.SetOutput(f)
	//log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.DebugLevel)

	// PROGRAM

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, 0755)
	}

	mapper := map[string]func(s string, c interface{}){
		Tr("csv_header_item_product"):                 func(s string, c interface{}) { c.(*ConfigProd).Prod = s },
		Tr("csv_header_stock_unit"):                   func(s string, c interface{}) { c.(*ConfigProd).StockUnit = s },
		Tr("csv_header_order_unit"):                   func(s string, c interface{}) { c.(*ConfigProd).OrderUnit = s },
		Tr("csv_header_factor_stock_unit_order_unit"): func(s string, c interface{}) { c.(*ConfigProd).Fact, _ = strconv.Atoi(s) },
		Tr("csv_header_stock_min"):                    func(s string, c interface{}) { c.(*ConfigProd).Min, _ = strconv.Atoi(s) },
	}

	config := make([]ConfigProd, 0)
	newLiner := func() interface{} { return new(ConfigProd) }
	appender := func(v interface{}) {
		a := v.(*ConfigProd)
		config = append(config, *a)
	}

	configFilename := configPrefix + Tr("file_name_config") + extension

	csvToStruct(configFilename, mapper, newLiner, appender)

	files, err3 := ioutil.ReadDir(dir)

	if err3 != nil {
		log.WithFields(log.Fields{
			"dir": dir,
			"err": err3,
		}).Error(Tr("error_dir_read"))
	}

	RouteFile(files)

	iStock := inStock(config).Copy()

	prodValHeader := []string{Tr("csv_header_item_product"), Tr("csv_header_item_value")}
	prodValRender := func(tab *strtab.T) [][]string { return tab.GetContentWithHeaders(false) }
	prodEvolRender := func(tab *strtab.T) [][]string { return tab.GetContentWithHeaders(true) }

	TableView{"bievre", generatedPrefix + Tr("file_name_stock") + extension,
		strtab.NewTable(prodValHeader, mapItem(iStock)...), prodValRender}.Show()

	TableView{"bievre", draftFileName(3, Tr("file_name_inventory")),
		strtab.NewTable(prodValHeader, mapItem(iStock)...), prodValRender}.Show()

	//Missing
	missing := missing(config, iStock)

	TableView{"bievre", generatedPrefix + Tr("file_name_to_order") + extension,
		strtab.NewTable(prodValHeader, mapItem(missing)...), prodValRender}.Show()

	TableView{"bievre", draftFileName(4, Tr("file_name_order")),
		strtab.NewTable(prodValHeader, mapItem(missing)...), prodValRender}.Show()

	TableView{"bievre", draftFileName(1, Tr("file_name_stock_in")),
		strtab.NewTable(prodValHeader, mapConfigProd(config)...), prodValRender}.Show()

	TableView{"bievre", draftFileName(2, Tr("file_name_stock_out")),
		strtab.NewTable(prodValHeader, mapConfigProd(config)...), prodValRender}.Show()

	TableView{"bievre", generatedPrefix + Tr("file_name_product") + extension,
		strtab.NewTableFromMap(mapItemsMap(endPt.ProdValEvolution("bievre"))).Transpose(), prodEvolRender}.Show()
}

func draftFileName(num int, name string) string {
	today := time.Now().Format(TimeFormat)
	return Tr("file_prefix_draft") + sep + today + sep + numberPrefix + strconv.Itoa(num) + sep + name + extension
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

		skelet.ExecuteCommand(skelet.Cmd{T: out, Route: stockRoute}, stock.Chain)
	}
}

func UnmarshalCsvFile(path Filename) (out skelet.Ider, err error) {

	// should put this in a local type
	mapper := map[string]func(s string, c interface{}){
		Tr("csv_header_item_product"): func(s string, c interface{}) { c.(*Item).Prod = Prod(s) },
		Tr("csv_header_item_value"): func(s string, c interface{}) {
			val, _ := strconv.Atoi(s)
			c.(*Item).Val = NewVal(UnitVal{base, val})
		},
	}

	its := make([]Item, 0)
	newLiner := func() interface{} { return new(Item) }
	appender := func(v interface{}) {
		a := v.(*Item)
		its = append(its, *a)
	}

	// should put this inside switch case
	csvToStruct(path.String(), mapper, newLiner, appender)
	itsMap := itemArrayToMap(its)

	switch path.Act {
	case Tr("file_name_stock_in"):
		return stockSk.InCmd{path.Stock, itsMap, path.Date.Format(TimeFormat) + sep + numberPrefix + path.Id}, nil
	case Tr("file_name_stock_out"):
		return stockSk.OutCmd{path.Stock, itsMap, path.Date.Format(TimeFormat) + sep + numberPrefix + path.Id}, nil
	case Tr("file_name_inventory"):
		return stockSk.InventoryCmd{path.Stock, itsMap, path.Date.Format(TimeFormat) + sep + numberPrefix + path.Id}, nil

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
	f.Stock = "bievre" //To refactor

	//Id
	f.Id = strings.TrimPrefix(args[3], numberPrefix)

	//Action
	f.Act = args[4]

	return
}

type ConfigProd struct {
	Prod      string
	StockUnit string
	OrderUnit string
	Fact      int
	Min       int
}

func GetMissingItems(c []ConfigProd) (out Items) {
	out = map[string]Item{}
	for _, p := range c {
		out[p.Prod] = Item{Prod(p.Prod), NewVal(UnitVal{base, p.Min * p.Fact})}
	}
	return
}

type Filename struct {
	Date                        time.Time
	Stock, Act, Id, Ext, Status string
}

func (f Filename) String() string {
	s := f.Date.Format(TimeFormat) + sep + numberPrefix + f.Id + sep + f.Act
	if f.Status != "" {
		s = s + sep + f.Status
	}
	s = s + f.Ext

	return s
}
