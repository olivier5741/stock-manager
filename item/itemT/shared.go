package itemT

import (
	"github.com/olivier5741/stock-manager/item/unitval"
	"github.com/olivier5741/stock-manager/item/val"
	"testing"
)

var (
	Inconnu  = unitval.Unit{"inconnu", 0}
	Pillule  = unitval.Unit{"pill.", 1}
	Tablette = unitval.Unit{"tab.", 15}
	Boite    = unitval.Unit{"b.", 45}
	Carton   = unitval.Unit{"cart.", 450}
)

func ValEqualCheck(t *testing.T, gots, exps val.Val) {

	for key, got := range gots.Vals {
		if exp, ok := exps.Vals[key]; !ok {
			t.Errorf("Unit %+v from val %+v does not exist in expected val %+v", got, gots, exps)
		} else {
			UnitValCheck(t, got, exp)
		}
	}

	for key, exp := range exps.Vals {
		if _, ok := gots.Vals[key]; !ok {
			t.Errorf("Unit %+v from expected val %+v does not exist in got val %+v", exp, exps, gots)
		}
	}

}

func UnitValCheck(t *testing.T, got, exp unitval.T) {
	if got.Fact != exp.Fact {
		t.Errorf("Fact %+v of unit %+v is not the same as expected %+v (unit %+v)", got.Fact, got, exp.Fact, exp)
	}

	if got.Val != exp.Val {
		t.Errorf("Val %+v of unit %+v is not the same as expected %+v (unit %+v)", got.Val, got, exp.Val, exp)
	}
}
