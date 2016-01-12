package order

import (
	. "github.com/olivier5741/stock-manager/item"
	. "github.com/olivier5741/stock-manager/item_test_shared"
	"testing"
)

func TestMissing(t *testing.T) {
	stock := Items{
		AspK: Asp5,
		IsoK: Iso3,
	}

	min := Items{
		AspK: Asp20,
		IsoK: Iso10,
	}

	gots := stock.Missing(min)

	exps := []Item{Asp15, Iso7}

	for _, exp := range exps {
		got := gots[string(exp.Prod)]
		if exp.Val != got.Val {
			CheckItemVal(t, got, exp)
		}
	}

}
