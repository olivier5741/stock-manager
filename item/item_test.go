package item

import (
	"testing"
)

var (
	Inconnu  = Unit{"inconnu", 0}
	Cachet   = Unit{"cachet", 1}
	Tablette = Unit{"Tablette", 15}
	Boite    = Unit{"boîte", 45}
	Carton   = Unit{"Carton", 450}
)

func TestAddVal(t *testing.T) {
	v1 := NewValFromUnitVals([]UnitVal{
		{Cachet, 12},
		{Tablette, 2},
		{Boite, 2},
		{Carton, 1},
	}...)

	v2 := NewValFromUnitVals([]UnitVal{
		{Tablette, 3},
		{Boite, 2},
	}...)

	got, _ := AddVal(v1, v2)

	exp := NewValFromUnitVals([]UnitVal{
		{Cachet, 12},
		{Tablette, 5},
		{Boite, 4},
		{Carton, 1},
	}...)

	ValEqualCheck(t, got, exp)
}

func TestSubVal(t *testing.T) {
	v1 := NewValFromUnitVals([]UnitVal{
		{Cachet, 12},
		{Tablette, 2},
		{Boite, 2},
		{Carton, 1},
	}...)

	v2 := NewValFromUnitVals([]UnitVal{
		{Cachet, 16},
		{Tablette, 3},
		{Boite, 2},
	}...)

	got, _ := SubVal(v1, v2)

	exp := NewValFromUnitVals([]UnitVal{
		{Cachet, 11},
		{Tablette, 1},
		{Boite, 9},
		{Carton, 0},
	}...)

	ValEqualCheck(t, got, exp)
}

func TestRedistribute(t *testing.T) {
	// l'unité principale ne peut valoir 0

	v1 := NewValFromUnitVals([]UnitVal{
		{Cachet, 50},
		{Tablette, 3},
		{Boite, 2},
	}...)

	got := v1.Redistribute()

	exp := NewValFromUnitVals([]UnitVal{
		{Cachet, 5},
		{Tablette, 0},
		{Boite, 4},
	}...)

	ValEqualCheck(t, got, exp)
}

func TestTotal(t *testing.T) {
	v1 := NewValFromUnitVals([]UnitVal{
		{Inconnu, 21},
		{Cachet, 50},
		{Tablette, 3},
		{Boite, 2},
	}...)

	got := v1.Total()

	exp := NewValFromUnitVals([]UnitVal{
		{Inconnu, 21},
		{Cachet, 185},
	}...)

	ValEqualCheck(t, got, exp)
}

func ValEqualCheck(t *testing.T, gots, exps Val) {
	if gots.Main != exps.Main {
		t.Errorf("Main %+v of val %+v is not the same as expected %+v (val %+v)", gots.Main, gots, exps.Main, exps)
	}

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

func UnitValCheck(t *testing.T, got, exp UnitVal) {
	if got.Fact != exp.Fact {
		t.Errorf("Fact %+v of unit %+v is not the same as expected %+v (unit %+v)", got.Fact, got, exp.Fact, exp)
	}

	if got.Val != exp.Val {
		t.Errorf("Val %+v of unit %+v is not the same as expected %+v (unit %+v)", got.Val, got, exp.Val, exp)
	}
}
