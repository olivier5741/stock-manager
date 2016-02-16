package stockT

import (
	stockCmd "github.com/olivier5741/stock-manager/cmd/stock"
	"github.com/olivier5741/stock-manager/item/itemT"
	"github.com/olivier5741/stock-manager/item/items"
	"github.com/olivier5741/stock-manager/item"
	"github.com/olivier5741/stock-manager/item/quant"
	"github.com/olivier5741/stock-manager/item/amount"
	sk "github.com/olivier5741/stock-manager/skelet"
	stockSk "github.com/olivier5741/stock-manager/stock/skelet"
	"testing"
)

var (
	AspK  = "aspirine"
	Asp   = item.Prod(AspK)
	Asp1  = item.Item{Asp, amount.NewA(quant.Q{itemT.Pillule, 1})}
	Asp5  = item.Item{Asp, amount.NewA(quant.Q{itemT.Pillule, 5})}
	Asp6  = item.Item{Asp, amount.NewA(quant.Q{itemT.Pillule, 6})}
	Asp8  = item.Item{Asp, amount.NewA(quant.Q{itemT.Pillule, 8})}
	Asp15 = item.Item{Asp, amount.NewA(quant.Q{itemT.Pillule, 15})}
	Asp20 = item.Item{Asp, amount.NewA(quant.Q{itemT.Pillule, 20})}
	IsoK  = "isobétadine"
	Iso   = item.Prod(IsoK)
	Iso0  = item.Item{Iso, amount.NewA(quant.Q{itemT.Pillule, 0})}
	Iso1  = item.Item{Iso, amount.NewA(quant.Q{itemT.Pillule, 1})}
	Iso2  = item.Item{Iso, amount.NewA(quant.Q{itemT.Pillule, 2})}
	Iso3  = item.Item{Iso, amount.NewA(quant.Q{itemT.Pillule, 3})}
	Iso4  = item.Item{Iso, amount.NewA(quant.Q{itemT.Pillule, 4})}
	Iso7  = item.Item{Iso, amount.NewA(quant.Q{itemT.Pillule, 7})}
	Iso10 = item.Item{Iso, amount.NewA(quant.Q{itemT.Pillule, 10})}
	Iso20 = item.Item{Iso, amount.NewA(quant.Q{itemT.Pillule, 20})}

	repo       = stockCmd.MakeDummyStockRepository()
	e          = stockCmd.EndPt{Db: repo}
	stockRoute = func(t sk.Ider) (ok bool, a sk.AggAct, p sk.EvtSrcPersister) {
		switch t.(type) {
		case stockSk.InCmd:
			return true, e.HandleIn, repo
		case stockSk.OutCmd:
			return true, e.HandleOut, repo
		case stockSk.InventoryCmd:
			return true, e.HandleInventory, repo
		default:
			return false, nil, nil
		}
	}
)

func CheckItemsValueAndExistence(t *testing.T, gots items.T, exps items.T, name string) {
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
