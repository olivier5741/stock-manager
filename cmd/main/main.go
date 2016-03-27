package main

import (
	log "github.com/Sirupsen/logrus"
	// considering using this one instead : i18n4go, tools better but not inside
	"os"
	"github.com/olivier5741/stock-manager/asset"
	stockCmd "github.com/olivier5741/stock-manager/cmd/stock"
	"github.com/olivier5741/stock-manager/item/items"
	"github.com/olivier5741/stock-manager/skelet"
	"github.com/olivier5741/stock-manager/stock"
	"github.com/olivier5741/stock-manager/port/sheet"
//	"github.com/olivier5741/stock-manager/port/sheet/osfile"
	"github.com/olivier5741/stock-manager/port/sheet/drivefile"
	"net/http"
	"github.com/codegangsta/negroni"
	"github.com/olivier5741/strtab"
)

// vim -u ~/.vimrc.go
// GOOS=windows GOARCH=386 go build -o stock-manager-0.1.exe c-main.go

var (

	//rickyFolder = "stock-ricky/"
	//mickyFolder = "commande-micky/"
	//dasdboardData = "données-dashboard/"

	//rickyAcquire = osfile.OsFile{rickyFolder}
	//rickyAnalyse = osfile.OsFile{rickyFolder+dasdboardData}
	//mickyAcquire = osfile.OsFile{mickyFolder}
	//mickyAnalyse = osfile.OsFile{mickyFolder+dasdboardData}

	rickyAcquire = drivefile.DriveFile{"0BzIZ3dfuz-CEQVg5YU1Ia0dMY3c", drivefile.GetService(), make(map[string]string)}
	rickyAnalyse = drivefile.DriveFile{"0BzIZ3dfuz-CEaGh4OHFhWVBHZHc", drivefile.GetService(), make(map[string]string)}
	mickyAcquire = drivefile.DriveFile{"0BzIZ3dfuz-CEU0lHYUdKTVJ5aXM", drivefile.GetService(), make(map[string]string)}

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
		case stock.ProdsUpdateCmd:
			return true, endPt.HandleProdsUpdate, repo
		default:
			return false, nil, nil
		}
	}

	stockId = "main"

	prodValHeader []string
)

func main() {
mux := http.NewServeMux()
  mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
    action()
  })

  n := negroni.Classic()
  n.UseHandler(mux)
  n.Run(":3000")
}

func action() {

	prodValHeader = []string{asset.Tr("csv_header_item_product"),
		asset.Tr("csv_header_item_value", 1), asset.Tr("csv_header_item_unit", 1),
		asset.Tr("csv_header_item_value", 2), asset.Tr("csv_header_item_unit", 2),
		asset.Tr("csv_header_item_value", 3), asset.Tr("csv_header_item_unit", 3),
		asset.Tr("csv_header_item_value", 4), asset.Tr("csv_header_item_unit", 4),
	}

	// LOGGING
	f, err1 := os.OpenFile(asset.Tr("file_name_log"),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err1 != nil {
		log.WithFields(log.Fields{
			"filename": asset.Tr("file_name_log"),
			"err":      err1,
		}).Error(asset.Tr("error_file_open"))
	}
	defer f.Close()

	log.SetOutput(f)
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.ErrorLevel)

	// UPDATE 
	all := sheet.AllSheets(rickyAcquire)

	for _,name := range all {

		var(
			cmd skelet.Ider
			its items.Items
			inN = asset.Tr("file_name_stock_in")
			outN = asset.Tr("file_name_stock_out")
			invN = asset.Tr("file_name_inventory")
			updateN = asset.Tr("file_name_product_update")
		)

		s := sheet.NewSheet(name,rickyAcquire)

		switch s.Name.Act {
		case inN:
			its = items.FromStringTable(s.Table.GetContentWithRowHeader())
			cmd = stock.InCmd{stockId, its, s.Name.Time()}
		case outN:
			its = items.FromStringTable(s.Table.GetContentWithRowHeader())
			cmd = stock.OutCmd{stockId, its, s.Name.Time()}
		case invN:
			its = items.FromStringTable(s.Table.GetContentWithRowHeader())
			cmd = stock.InventoryCmd{stockId, its, s.Name.Time()}
		case updateN:
			mins,units := stock.ProdsUpdateFromStringTable(s.Table.GetContentWithRowHeader())
			cmd = stock.ProdsUpdateCmd{stockId, mins, units, s.Name.Time()}
		default:
			log.Error(asset.Tr("no_action_for_filename_error"))
			continue
		}
		skelet.ExecuteCommand(skelet.Cmd{T: cmd, Route: stockRoute}, stockCmd.Chain)
	}


	// POPULATE VIEW
	stockInt, err4 := endPt.Db.Get("main")
	iStock := stockInt.(*stock.Stock).Items
	min := stockInt.(*stock.Stock).Min
	if err4 != nil {
		log.WithFields(log.Fields{
			"err": err4,
		}).Error(asset.Tr("error_query_stock"))
	}

	prodValRender := func(tab *strtab.T) [][]string {
		return tab.GetContentWithHeaders(false)
	}
	prodEvolRender := func(tab *strtab.T) [][]string {
		return tab.GetContentWithHeaders(true)
	}

	// not very good
	sheet.AllSheets(rickyAnalyse)
	sheet.AllSheets(mickyAcquire)

	sheet.Sheet{
		sheet.NewBasicFilename(asset.Tr("file_name_stock")),
		strtab.NewT(prodValHeader, iStock.StringSlice()...).Sort(),
		prodValRender}.Put(rickyAnalyse)

	sheet.Sheet{
		sheet.NewBasicFilename(asset.Tr("file_name_to_order")),
		strtab.NewT(prodValHeader, items.Missing(iStock,min).StringSlice()...).Sort(),
		prodValRender}.Put(mickyAcquire)

	sheet.Sheet{
		sheet.NewDraftFilename(3, asset.Tr("file_name_inventory")),
		strtab.NewT(prodValHeader, iStock.StringSlice()...).Sort(),
		prodValRender}.Put(rickyAcquire)

	sheet.Sheet{
		sheet.NewDraftFilename(1, asset.Tr("file_name_stock_in")),
		strtab.NewT(prodValHeader, iStock.Empty().StringSlice()...).Sort(),
		prodValRender}.Put(rickyAcquire)

	sheet.Sheet{
		sheet.NewDraftFilename(2, asset.Tr("file_name_stock_out")),
		strtab.NewT(prodValHeader, iStock.Empty().StringSlice()...).Sort(),
		prodValRender}.Put(rickyAcquire)

	sheet.Sheet{
		sheet.NewDraftFilename(4, asset.Tr("file_name_product_update")),
		strtab.NewT(prodValHeader, min.Empty().StringSlice()...).Sort(),
		prodValRender}.Put(rickyAcquire)

	sheet.Sheet{
		sheet.NewBasicFilename(asset.Tr("file_name_product_evolution")),
		strtab.NewTfromMap(items.ItemsMapToStringMapTable(
			endPt.ProdValEvol("main"))).Sort().Transpose().Sort(),
		prodEvolRender}.Put(rickyAnalyse)
}

