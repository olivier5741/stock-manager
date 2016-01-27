package itemtest

import (
	. "github.com/olivier5741/stock-manager/item"
	"testing"
)

func TestAddVal(t *testing.T) {
	v1 := NewVal([]UnitVal{
		{Pillule, 12},
		{Tablette, 2},
		{Boite, 2},
		{Carton, 1},
	}...)

	t.Log(v1.String())

	v2 := NewVal([]UnitVal{
		{Tablette, 3},
		{Boite, 2},
	}...)

	t.Log(v2.String())

	got := AddVal(v1, v2)

	t.Log(got.String())
	t.Log(got.TotalWithRound(Carton))

	exp := NewVal([]UnitVal{
		{Pillule, 12},
		{Tablette, 5},
		{Boite, 4},
		{Carton, 1},
	}...)

	ValEqualCheck(t, got, exp)
}

func TestSubVal(t *testing.T) {
	v1 := NewVal([]UnitVal{
		{Pillule, 12},
		{Tablette, 2},
		{Boite, 2},
		{Carton, 1},
	}...)

	t.Log(v1.String())

	v2 := NewVal([]UnitVal{
		{Pillule, 16},
		{Tablette, 3},
		{Boite, 2},
	}...)

	t.Log(v2.String())

	got := SubVal(v1, v2)

	t.Log(got.String())
	t.Log(got.TotalWithRound(Carton))

	exp := NewVal([]UnitVal{
		{Pillule, 11},
		{Tablette, 1},
		{Boite, 9},
		{Carton, 0},
	}...)

	ValEqualCheck(t, got, exp)
}

func TestDiff(t *testing.T) {
	v1 := NewVal([]UnitVal{
		{Pillule, 12},
		{Tablette, 40},
		{Boite, 2},
		{Carton, 1},
	}...)

	t.Log(v1.String())

	v2 := NewVal([]UnitVal{
		{Pillule, 16},
		{Tablette, 3},
		{Boite, 2},
	}...)

	t.Log(v2.String())

	got, _, _ := Diff(v1, v2)

	t.Log(got.String())
	t.Log(got.TotalWithRound(Carton))

	exp := NewVal([]UnitVal{
		{Carton, 2},
		{Boite, 2},
		{Tablette, 0},
		{Pillule, 11},
	}...)

	ValEqualCheck(t, got, exp)
}

func TestRedistribute(t *testing.T) {
	// l'unité principale ne peut valoir 0

	v1 := NewVal([]UnitVal{
		{Pillule, 50},
		{Tablette, 3},
		{Boite, 2},
	}...)

	t.Log(v1.String())

	got := v1.Redistribute()

	t.Log(got.String())
	t.Log(got.TotalWithRound(Boite))

	exp := NewVal([]UnitVal{
		{Pillule, 5},
		{Tablette, 0},
		{Boite, 4},
	}...)

	ValEqualCheck(t, got, exp)
}

func TestTotal(t *testing.T) {
	v1 := NewVal([]UnitVal{
		{Inconnu, 21},
		{Pillule, 50},
		{Tablette, 3},
		{Boite, 2},
	}...)

	t.Log(v1.String())

	got := v1.Total()

	t.Log(got.String())
	t.Log(got.TotalWithRound(Boite))

	exp := NewVal([]UnitVal{
		{Inconnu, 21},
		{Pillule, 185},
	}...)

	ValEqualCheck(t, got, exp)
}
