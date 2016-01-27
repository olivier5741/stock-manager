package stocktest

import (
	. "github.com/olivier5741/stock-manager/cmd/stock"
	. "github.com/olivier5741/stock-manager/item"
	. "github.com/olivier5741/stock-manager/skelet"
	. "github.com/olivier5741/stock-manager/stock/main"
	. "github.com/olivier5741/stock-manager/stock/skelet"
	"testing"
)

func TestMain(t *testing.T) {

	// It could be interesting to retain 'Route' in some kind of object

	cmd1 := Cmd{
		T:     InCmd{"Carlsbourg", Items{IsoK: Iso4}, "2016-01-27"},
		Route: stockRoute,
	}

	cmd2 := Cmd{
		T:     OutCmd{"Carlsbourg", Items{IsoK: Iso1}, "2016-01-28"},
		Route: stockRoute,
	}

	cmd3 := Cmd{
		T:     InventoryCmd{"Carlsbourg", Items{IsoK: Iso2}, "2016-01-29"},
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
