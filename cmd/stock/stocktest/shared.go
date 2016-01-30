package stocktest

import (
	. "github.com/olivier5741/stock-manager/cmd/stock"
	. "github.com/olivier5741/stock-manager/item"
	. "github.com/olivier5741/stock-manager/item/itemtest"
	. "github.com/olivier5741/stock-manager/skelet"
	. "github.com/olivier5741/stock-manager/stock/skelet"
	"testing"
)

var (
	AspK  = "aspirine"
	Asp   = Prod(AspK)
	Asp1  = Item{Asp, NewVal(UnitVal{Pillule, 1})}
	Asp5  = Item{Asp, NewVal(UnitVal{Pillule, 5})}
	Asp6  = Item{Asp, NewVal(UnitVal{Pillule, 6})}
	Asp8  = Item{Asp, NewVal(UnitVal{Pillule, 8})}
	Asp15 = Item{Asp, NewVal(UnitVal{Pillule, 15})}
	Asp20 = Item{Asp, NewVal(UnitVal{Pillule, 20})}
	IsoK  = "isobétadine"
	Iso   = Prod(IsoK)
	Iso0  = Item{Iso, NewVal(UnitVal{Pillule, 0})}
	Iso1  = Item{Iso, NewVal(UnitVal{Pillule, 1})}
	Iso2  = Item{Iso, NewVal(UnitVal{Pillule, 2})}
	Iso3  = Item{Iso, NewVal(UnitVal{Pillule, 3})}
	Iso4  = Item{Iso, NewVal(UnitVal{Pillule, 4})}
	Iso7  = Item{Iso, NewVal(UnitVal{Pillule, 7})}
	Iso10 = Item{Iso, NewVal(UnitVal{Pillule, 10})}
	Iso20 = Item{Iso, NewVal(UnitVal{Pillule, 20})}

	repo       = MakeDummyStockRepository()
	e          = EndPt{Db: repo}
	stockRoute = func(t Ider) (ok bool, a AggAct, p EvtSrcPersister) {
		switch t.(type) {
		case InCmd:
			return true, e.HandleIn, repo
		case OutCmd:
			return true, e.HandleOut, repo
		case InventoryCmd:
			return true, e.HandleInventory, repo
		default:
			return false, nil, nil
		}
	}
)

func CheckItemsValueAndExistence(t *testing.T, gots, exps Items, name string) {
	for _, exp := range exps {
		got, ok := gots[string(exp.Prod)]
		if !ok {
			t.Errorf("%q does not exist in "+name, exp.Prod)
		}
		if _, no, diff := Diff(got.Val, exp.Val); !(no && diff == 0) {
			t.Errorf("not equal")
		}

	}
}
