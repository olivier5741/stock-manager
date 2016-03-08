package itemT

import (
	"github.com/olivier5741/stock-manager/item/amount"
	"github.com/olivier5741/stock-manager/item/quant"
	"testing"
	"math/big"
)

var (
	Inconnu  = quant.Unit{"inconnu", new(big.Rat).SetInt64(0)}
	Pillule  = quant.Unit{"pill.", new(big.Rat).SetInt64(1)}
	Tablette = quant.Unit{"tab.", new(big.Rat).SetInt64(15)}
	Boite    = quant.Unit{"b.", new(big.Rat).SetInt64(45)}
	Carton   = quant.Unit{"cart.", new(big.Rat).SetInt64(450)}
)

func ValEqualCheck(t *testing.T, gots, exps amount.Amount) {

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

func UnitValCheck(t *testing.T, got, exp quant.Quant) {
	if got.Fact.Cmp(exp.Fact) != 0 {
		t.Errorf("Fact %+v of unit %+v is not the same as expected %+v (unit %+v)", got.Fact, got, exp.Fact, exp)
	}

	if got.Val.Cmp(exp.Val) != 0 {
		t.Errorf("Val %+v of unit %+v is not the same as expected %+v (unit %+v)", got.Val, got, exp.Val, exp)
	}
}
