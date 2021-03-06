package stockT

import (
	stockCmd "github.com/olivier5741/stock-manager/cmd/stock"
	"github.com/olivier5741/stock-manager/item"
	"github.com/olivier5741/stock-manager/item/amount"
	"github.com/olivier5741/stock-manager/item/itemT"
	"github.com/olivier5741/stock-manager/item/items"
	"github.com/olivier5741/stock-manager/item/quant"
	"github.com/olivier5741/stock-manager/skelet"
	"github.com/olivier5741/stock-manager/stock"
	"testing"
)

var (
	AspK  = "aspirine"
	Asp   = item.Prod(AspK)
	Asp1  = item.Item{Asp, amount.NewAmount(quant.Quant{itemT.Pillule, 1})}
	Asp5  = item.Item{Asp, amount.NewAmount(quant.Quant{itemT.Pillule, 5})}
	Asp6  = item.Item{Asp, amount.NewAmount(quant.Quant{itemT.Pillule, 6})}
	Asp8  = item.Item{Asp, amount.NewAmount(quant.Quant{itemT.Pillule, 8})}
	Asp15 = item.Item{Asp, amount.NewAmount(quant.Quant{itemT.Pillule, 15})}
	Asp20 = item.Item{Asp, amount.NewAmount(quant.Quant{itemT.Pillule, 20})}
	IsoK  = "isobétadine"
	Iso   = item.Prod(IsoK)
	Iso0  = item.Item{Iso, amount.NewAmount(quant.Quant{itemT.Pillule, 0})}
	Iso1  = item.Item{Iso, amount.NewAmount(quant.Quant{itemT.Pillule, 1})}
	Iso2  = item.Item{Iso, amount.NewAmount(quant.Quant{itemT.Pillule, 2})}
	Iso3  = item.Item{Iso, amount.NewAmount(quant.Quant{itemT.Pillule, 3})}
	Iso4  = item.Item{Iso, amount.NewAmount(quant.Quant{itemT.Pillule, 4})}
	Iso7  = item.Item{Iso, amount.NewAmount(quant.Quant{itemT.Pillule, 7})}
	Iso10 = item.Item{Iso, amount.NewAmount(quant.Quant{itemT.Pillule, 10})}
	Iso20 = item.Item{Iso, amount.NewAmount(quant.Quant{itemT.Pillule, 20})}

	repo       = stockCmd.MakeDummyStockRepo()
	e          = stockCmd.EndPt{Db: repo}
	stockRoute = func(t skelet.Ider) (ok bool, a skelet.AggAct, p skelet.EvtSrcPersister) {
		switch t.(type) {
		case stock.InCmd:
			return true, e.HandleIn, repo
		case stock.OutCmd:
			return true, e.HandleOut, repo
		case stock.InventoryCmd:
			return true, e.HandleInventory, repo
		default:
			return false, nil, nil
		}
	}
)

func CheckItemsValueAndExistence(t *testing.T, gots items.Items, exps items.Items, name string) {
	for _, exp := range exps {
		got, ok := gots[string(exp.Prod)]
		if !ok {
			t.Errorf("%q does not exist in "+name, exp.Prod)
		}
		if _, no, diff := amount.Diff(got.Amount, exp.Amount); !(no && diff == 0) {
			t.Errorf("not equal")
		}

	}
}
