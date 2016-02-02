package stockT

import (
	stockCmd "github.com/olivier5741/stock-manager/cmd/stock"
	"github.com/olivier5741/stock-manager/item/items"
	sk "github.com/olivier5741/stock-manager/skelet"
	stock "github.com/olivier5741/stock-manager/stock/main"
	stockSk "github.com/olivier5741/stock-manager/stock/skelet"
	"testing"
)

func TestMain(t *testing.T) {

	// It could be interesting to retain 'Route' in some kind of object

	cmd1 := sk.Cmd{
		T:     stockSk.InCmd{"Carlsbourg", items.T{IsoK: Iso4}, "2016-01-27"},
		Route: stockRoute,
	}

	cmd2 := sk.Cmd{
		T:     stockSk.OutCmd{"Carlsbourg", items.T{IsoK: Iso1}, "2016-01-28"},
		Route: stockRoute,
	}

	cmd3 := sk.Cmd{
		T:     stockSk.InventoryCmd{"Carlsbourg", items.T{IsoK: Iso2}, "2016-01-29"},
		Route: stockRoute,
	}

	sk.ExecuteCommand(cmd1, stockCmd.Chain)
	sk.ExecuteCommand(cmd2, stockCmd.Chain)
	sk.ExecuteCommand(cmd3, stockCmd.Chain)

	c, err := repo.Get("Carlsbourg")

	if err != nil {
		t.Fatal(err)
	}

	s := c.(*stock.Stock)

	exps := map[string]items.Item{IsoK: Iso2}

	CheckItemsValueAndExistence(t, s.T, exps, "stock")
}
