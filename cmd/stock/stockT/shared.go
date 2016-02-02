package stockT

import (
	stockCmd "github.com/olivier5741/stock-manager/cmd/stock"
	"github.com/olivier5741/stock-manager/item/itemT"
	"github.com/olivier5741/stock-manager/item/items"
	"github.com/olivier5741/stock-manager/item/unitval"
	"github.com/olivier5741/stock-manager/item/val"
	sk "github.com/olivier5741/stock-manager/skelet"
	stockSk "github.com/olivier5741/stock-manager/stock/skelet"
	"testing"
)

var (
	AspK  = "aspirine"
	Asp   = items.Prod(AspK)
	Asp1  = items.Item{Asp, val.NewT(unitval.T{itemT.Pillule, 1})}
	Asp5  = items.Item{Asp, val.NewT(unitval.T{itemT.Pillule, 5})}
	Asp6  = items.Item{Asp, val.NewT(unitval.T{itemT.Pillule, 6})}
	Asp8  = items.Item{Asp, val.NewT(unitval.T{itemT.Pillule, 8})}
	Asp15 = items.Item{Asp, val.NewT(unitval.T{itemT.Pillule, 15})}
	Asp20 = items.Item{Asp, val.NewT(unitval.T{itemT.Pillule, 20})}
	IsoK  = "isobétadine"
	Iso   = items.Prod(IsoK)
	Iso0  = items.Item{Iso, val.NewT(unitval.T{itemT.Pillule, 0})}
	Iso1  = items.Item{Iso, val.NewT(unitval.T{itemT.Pillule, 1})}
	Iso2  = items.Item{Iso, val.NewT(unitval.T{itemT.Pillule, 2})}
	Iso3  = items.Item{Iso, val.NewT(unitval.T{itemT.Pillule, 3})}
	Iso4  = items.Item{Iso, val.NewT(unitval.T{itemT.Pillule, 4})}
	Iso7  = items.Item{Iso, val.NewT(unitval.T{itemT.Pillule, 7})}
	Iso10 = items.Item{Iso, val.NewT(unitval.T{itemT.Pillule, 10})}
	Iso20 = items.Item{Iso, val.NewT(unitval.T{itemT.Pillule, 20})}

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
		if _, no, diff := val.Diff(got.Val, exp.Val); !(no && diff == 0) {
			t.Errorf("not equal")
		}

	}
}
