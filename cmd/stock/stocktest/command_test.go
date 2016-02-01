
package stocktest

import (
	item "github.com/olivier5741/stock-manager/item"
	sk "github.com/olivier5741/stock-manager/skelet"
	stock "github.com/olivier5741/stock-manager/stock/main"
	stockCmd "github.com/olivier5741/stock-manager/cmd/stock"
	stockSk "github.com/olivier5741/stock-manager/stock/skelet"
	"testing"
)

func TestMain(t *testing.T) {

	// It could be interesting to retain 'Route' in some kind of object

	cmd1 := sk.Cmd{
		T:     stockSk.InCmd{"Carlsbourg", item.Items{IsoK: Iso4}, "2016-01-27"},
		Route: stockRoute,
	}

	cmd2 := sk.Cmd{
		T:     stockSk.OutCmd{"Carlsbourg", item.Items{IsoK: Iso1}, "2016-01-28"},
		Route: stockRoute,
	}

	cmd3 := sk.Cmd{
		T:     stockSk.InventoryCmd{"Carlsbourg", item.Items{IsoK: Iso2}, "2016-01-29"},
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

	exps := map[string]item.Item{IsoK: Iso2}

	CheckItemsValueAndExistence(t, s.Items, exps, "stock")
}
