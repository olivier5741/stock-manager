package main

import (
	log "github.com/Sirupsen/logrus"
	// considering using this one instead : i18n4go, tools better but not inside
	"os"
	"github.com/nicksnyder/go-i18n/i18n"
	stockCmd "github.com/olivier5741/stock-manager/cmd/stock"
	"github.com/olivier5741/stock-manager/item/items"
	"github.com/olivier5741/stock-manager/skelet"
	"github.com/olivier5741/stock-manager/stock"
	"github.com/olivier5741/stock-manager/port/sheet"
	"github.com/olivier5741/stock-manager/port/sheet/osfile"
	"github.com/olivier5741/strtab"
	"github.com/olivier5741/stock-manager/tr"
)

// vim -u ~/.vimrc.go
// GOOS=windows GOARCH=386 go build -o stock-manager-0.1.exe c-main.go

var (
	Tr i18n.TranslateFunc

	loggingPrefix  = "l-"

	files = osfile.OsFile{"./"}

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

	stockId = "main"

	prodValHeader []string
)

func init() {
	prodValHeader = []string{tr.Tr("csv_header_item_product"),
		tr.Tr("csv_header_item_value", 1), tr.Tr("csv_header_item_unit", 1),
		tr.Tr("csv_header_item_value", 2), tr.Tr("csv_header_item_unit", 2),
		tr.Tr("csv_header_item_value", 3), tr.Tr("csv_header_item_unit", 3),
		tr.Tr("csv_header_item_value", 4), tr.Tr("csv_header_item_unit", 4),
	}
}

func main() {

	// LOGGING
	logfile := loggingPrefix + tr.Tr("file_name_log")
	f, err1 := os.OpenFile(logfile,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err1 != nil {
		log.WithFields(log.Fields{
			"filename": logfile,
			"err":      err1,
		}).Error(tr.Tr("error_file_open"))
	}
	defer f.Close()

	log.SetOutput(f)
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.ErrorLevel)

	// UPDATE 
	all := sheet.AllSheets(files)

	for _,name := range all {
		s := sheet.NewSheet(name,files)
		its := items.FromStringTable(s.Table.GetContentWithRowHeader())
		var cmd skelet.Ider

		switch s.Name.Act {
		case tr.Tr("file_name_stock_in"):
			cmd = stock.InCmd{stockId, its, s.Name.Time()}
		case tr.Tr("file_name_stock_out"):
			cmd = stock.OutCmd{stockId, its, s.Name.Time()}
		case tr.Tr("file_name_inventory"):
			cmd = stock.InventoryCmd{stockId, its, s.Name.Time()}
		default:
			log.Error(tr.Tr("no_action_for_filename_error"))
		}

		skelet.ExecuteCommand(skelet.Cmd{T: cmd, Route: stockRoute}, stockCmd.Chain)
	}


	// POPULATE VIEW
	stockInt, err4 := endPt.Db.Get("main")
	iStock := stockInt.(*stock.Stock).Items
	if err4 != nil {
		log.WithFields(log.Fields{
			"err": err4,
		}).Error(tr.Tr("error_query_stock"))
	}

	prodValRender := func(tab *strtab.T) [][]string {
		return tab.GetContentWithHeaders(false)
	}
	prodEvolRender := func(tab *strtab.T) [][]string {
		return tab.GetContentWithHeaders(true)
	}

	sheet.Sheet{
		sheet.NewBasicFilename(tr.Tr("file_name_stock")),
		strtab.NewT(prodValHeader, iStock.StringSlice()...).Sort(),
		prodValRender}.Put(files)

	sheet.Sheet{
		sheet.NewDraftFilename(3, tr.Tr("file_name_inventory")),
		strtab.NewT(prodValHeader, iStock.StringSlice()...).Sort(),
		prodValRender}.Put(files)

	sheet.Sheet{
		sheet.NewDraftFilename(1, tr.Tr("file_name_stock_in")),
		strtab.NewT(prodValHeader, iStock.Empty().StringSlice()...).Sort(),
		prodValRender}.Put(files)

	sheet.Sheet{
		sheet.NewDraftFilename(2, tr.Tr("file_name_stock_out")),
		strtab.NewT(prodValHeader, iStock.Empty().StringSlice()...).Sort(),
		prodValRender}.Put(files)

	sheet.Sheet{
		sheet.NewBasicFilename(tr.Tr("file_name_product")),
		strtab.NewTfromMap(items.ItemsMapToStringMapTable(
			endPt.ProdValEvol("main"))).Sort().Transpose().Sort(),
		prodEvolRender}.Put(files)
}

