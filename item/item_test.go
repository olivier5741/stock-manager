package item

import (
	"testing"
)

var (
	Inconnu  = Unit{"inconnu", 0}
	Cachet   = Unit{"pill.", 1}
	Tablette = Unit{"tab.", 15}
	Boite    = Unit{"b.", 45}
	Carton   = Unit{"cart.", 450}
)

func TestAddVal(t *testing.T) {
	v1 := NewValFromUnitVals([]UnitVal{
		{Cachet, 12},
		{Tablette, 2},
		{Boite, 2},
		{Carton, 1},
	}...)

	t.Log(v1.String())

	v2 := NewValFromUnitVals([]UnitVal{
		{Tablette, 3},
		{Boite, 2},
	}...)

	t.Log(v2.String())

	got, _ := AddVal(v1, v2)

	t.Log(got.String())

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

	t.Log(v1.String())

	v2 := NewValFromUnitVals([]UnitVal{
		{Cachet, 16},
		{Tablette, 3},
		{Boite, 2},
	}...)

	t.Log(v2.String())

	got, _ := SubVal(v1, v2)

	t.Log(got.String())

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

	t.Log(v1.String())

	got := v1.Redistribute()

	t.Log(got.String())

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

	t.Log(v1.String())

	got := v1.Total()

	t.Log(got.String())

	exp := NewValFromUnitVals([]UnitVal{
		{Inconnu, 21},
		{Cachet, 185},
	}...)

	ValEqualCheck(t, got, exp)
}

func ValEqualCheck(t *testing.T, gots, exps Val) {

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
