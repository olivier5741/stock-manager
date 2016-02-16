package stockT

import (
	stockCmd "github.com/olivier5741/stock-manager/cmd/stock"
	"github.com/olivier5741/stock-manager/item"
	"github.com/olivier5741/stock-manager/item/items"
	"github.com/olivier5741/stock-manager/skelet"
	"github.com/olivier5741/stock-manager/stock"
	"testing"
)

func TestMain(t *testing.T) {

	// It could be interesting to retain 'Route' in some kind of object

	cmd1 := skelet.Cmd{
		T:     stock.InCmd{"Carlsbourg", items.I{IsoK: Iso4}, "2016-01-27"},
		Route: stockRoute,
	}

	cmd2 := skelet.Cmd{
		T:     stock.OutCmd{"Carlsbourg", items.I{IsoK: Iso1}, "2016-01-28"},
		Route: stockRoute,
	}

	cmd3 := skelet.Cmd{
		T:     stock.InventoryCmd{"Carlsbourg", items.I{IsoK: Iso2}, "2016-01-29"},
		Route: stockRoute,
	}

	skelet.ExecuteCommand(cmd1, stockCmd.Chain)
	skelet.ExecuteCommand(cmd2, stockCmd.Chain)
	skelet.ExecuteCommand(cmd3, stockCmd.Chain)

	c, err := repo.Get("Carlsbourg")

	if err != nil {
		t.Fatal(err)
	}

	s := c.(*stock.Stock)

	exps := map[string]item.I{IsoK: Iso2}

	CheckItemsValueAndExistence(t, s.I, exps, "stock")
}
