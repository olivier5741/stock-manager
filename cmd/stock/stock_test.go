package stock

import (
	. "github.com/olivier5741/stock-manager/item"
	. "github.com/olivier5741/stock-manager/item_test_shared"
	. "github.com/olivier5741/stock-manager/skelet"
	. "github.com/olivier5741/stock-manager/stock"
	"testing"
)

func TestStockFromEvent(t *testing.T) {

	db := DummyStockRepository{db: map[string][]interface{}{}}

	stockEndPt := EndPt{db}

	stockEndPt.HandleIn(InCmd{"Bièvre", Items{AspK: Asp20, IsoK: Iso10}})

	db.Put("Bièvre", Out{
		Items: Items{
			AspK: Asp5,
			IsoK: Iso3},
	})

	db.Put("Bièvre", Out{
		Items: Items{
			AspK: Asp5,
			IsoK: Iso3},
	})

	db.Put("Bièvre", Inventory{
		Items: Items{
			AspK: Asp6,
			IsoK: Iso3},
	})

	db.Put("Bièvre", Out{
		Items: Items{
			AspK: Asp5,
			IsoK: Iso3},
	})

	s := db.Get("Bièvre").(Stock)

	exps := map[string]Item{AspK: Asp1, IsoK: Iso0}

	CheckItemsValueAndExistence(t, s.Items, exps, "stock")

}
