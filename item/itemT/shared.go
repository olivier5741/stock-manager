package itemT

import (
	"github.com/olivier5741/stock-manager/item/amount"
	"github.com/olivier5741/stock-manager/item/quant"
	"testing"
)

var (
	Inconnu  = quant.Unit{"inconnu", 0}
	Pillule  = quant.Unit{"pill.", 1}
	Tablette = quant.Unit{"tab.", 15}
	Boite    = quant.Unit{"b.", 45}
	Carton   = quant.Unit{"cart.", 450}
)

func ValEqualCheck(t *testing.T, gots, exps amount.A) {

	for key, got := range gots.QuantsMap() {
		if exp, ok := exps.QuantsMap()[key]; !ok {
			t.Errorf("Unit %+v from val %+v does not exist in expected val %+v", got, gots, exps)
		} else {
			UnitValCheck(t, got, exp)
		}
	}

	for key, exp := range exps.QuantsMap() {
		if _, ok := gots.QuantsMap()[key]; !ok {
			t.Errorf("Unit %+v from expected val %+v does not exist in got val %+v", exp, exps, gots)
		}
	}

}

func UnitValCheck(t *testing.T, got, exp quant.Q) {
	if got.Fact != exp.Fact {
		t.Errorf("Fact %+v of unit %+v is not the same as expected %+v (unit %+v)", got.Fact, got, exp.Fact, exp)
	}

	if got.Val != exp.Val {
		t.Errorf("Val %+v of unit %+v is not the same as expected %+v (unit %+v)", got.Val, got, exp.Val, exp)
	}
}
