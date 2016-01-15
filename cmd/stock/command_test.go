package stock

import (
	. "github.com/olivier5741/stock-manager/item"
	. "github.com/olivier5741/stock-manager/item_test_shared"
	. "github.com/olivier5741/stock-manager/skelet"
	. "github.com/olivier5741/stock-manager/stock"
	"testing"
)

var (
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

func TestMain(t *testing.T) {

	// It could be interesting to retain 'Route' in some kind of object

	cmd1 := Cmd{
		T:     InCmd{"Carlsbourg", Items{IsoK: Iso4}},
		Route: stockRoute,
	}

	cmd2 := Cmd{
		T:     OutCmd{"Carlsbourg", Items{IsoK: Iso1}},
		Route: stockRoute,
	}

	cmd3 := Cmd{
		T:     InventoryCmd{"Carlsbourg", Items{IsoK: Iso2}},
		Route: stockRoute,
	}

	ExecuteCommand(cmd1, Chain)
	ExecuteCommand(cmd2, Chain)
	ExecuteCommand(cmd3, Chain)

	c, err := repo.Get("Carlsbourg")

	if err != nil {
		t.Fatal(err)
	}

	s := c.(*Stock)

	exps := map[string]Item{IsoK: Iso2}

	CheckItemsValueAndExistence(t, s.Items, exps, "stock")
}
