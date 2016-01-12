package item_test_shared

import (
	. "github.com/olivier5741/stock-manager/item"
	"testing"
)

var (
	AspK  = "aspirine"
	Asp   = Prod(AspK)
	Asp1  = Item{Asp, 1}
	Asp5  = Item{Asp, 5}
	Asp6  = Item{Asp, 6}
	Asp8  = Item{Asp, 8}
	Asp15 = Item{Asp, 15}
	Asp20 = Item{Asp, 20}
	IsoK  = "isobétadine"
	Iso   = Prod(IsoK)
	Iso0  = Item{Iso, 0}
	Iso3  = Item{Iso, 3}
	Iso7  = Item{Iso, 7}
	Iso10 = Item{Iso, 10}
	Iso20 = Item{Iso, 20}
)

func CheckItemVal(t *testing.T, got, ex Item) {
	t.Errorf("Expected value for %q : %v not %v", ex.Prod, ex.Val, got.Val)
}

func CheckItemsValueAndExistence(t *testing.T, gots, exps Items, name string) {
	for _, exp := range exps {
		got, ok := gots[string(exp.Prod)]
		if !ok {
			t.Errorf("%q does not exist in "+name, exp.Prod)
		}
		if got.Val != exp.Val {
			CheckItemVal(t, got, exp)
		}

	}
}
