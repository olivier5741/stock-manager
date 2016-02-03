package itemT

import (
	"github.com/olivier5741/stock-manager/item/quant"
	"github.com/olivier5741/stock-manager/item/amount"
	"testing"
)

func TestAddVal(t *testing.T) {
	v1 := amount.NewA([]quant.Q{
		{Pillule, 12},
		{Tablette, 2},
		{Boite, 2},
		{Carton, 1},
	}...)

	t.Log(v1.String())

	v2 := amount.NewA([]quant.Q{
		{Tablette, 3},
		{Boite, 2},
	}...)

	t.Log(v2.String())

	got := amount.Add(v1, v2)

	t.Log(got.String())
	t.Log(got.TotalWithRound(Carton))

	exp := amount.NewA([]quant.Q{
		{Pillule, 12},
		{Tablette, 5},
		{Boite, 4},
		{Carton, 1},
	}...)

	ValEqualCheck(t, got, exp)
}

func TestSub(t *testing.T) {
	v1 := amount.NewA([]quant.Q{
		{Pillule, 12},
		{Tablette, 2},
		{Boite, 2},
		{Carton, 1},
	}...)

	t.Log(v1.String())

	v2 := amount.NewA([]quant.Q{
		{Pillule, 16},
		{Tablette, 3},
		{Boite, 2},
	}...)

	t.Log(v2.String())

	got := amount.Sub(v1, v2)

	t.Log(got.String())
	t.Log(got.TotalWithRound(Carton))

	exp := amount.NewA([]quant.Q{
		{Pillule, 11},
		{Tablette, 1},
		{Boite, 9},
		{Carton, 0},
	}...)

	ValEqualCheck(t, got, exp)
}

func TestDiff(t *testing.T) {
	v1 := amount.NewA([]quant.Q{
		{Pillule, 12},
		{Tablette, 40},
		{Boite, 2},
		{Carton, 1},
	}...)

	t.Log(v1.String())

	v2 := amount.NewA([]quant.Q{
		{Pillule, 16},
		{Tablette, 3},
		{Boite, 2},
	}...)

	t.Log(v2.String())

	got, _, _ := amount.Diff(v1, v2)

	t.Log(got.String())
	t.Log(got.TotalWithRound(Carton))

	exp := amount.NewA([]quant.Q{
		{Carton, 2},
		{Boite, 2},
		{Tablette, 0},
		{Pillule, 11},
	}...)

	ValEqualCheck(t, got, exp)
}

func TestRedistribute(t *testing.T) {
	// l'unité principale ne peut valoir 0

	v1 := amount.NewA([]quant.Q{
		{Pillule, 50},
		{Tablette, 3},
		{Boite, 2},
	}...)

	t.Log(v1.String())

	got := v1.Redistribute()

	t.Log(got.String())
	t.Log(got.TotalWithRound(Boite))

	exp := amount.NewA([]quant.Q{
		{Pillule, 5},
		{Tablette, 0},
		{Boite, 4},
	}...)

	ValEqualCheck(t, got, exp)
}

func TestTotal(t *testing.T) {
	v1 := amount.NewA([]quant.Q{
		{Inconnu, 21},
		{Pillule, 50},
		{Tablette, 3},
		{Boite, 2},
	}...)

	t.Log(v1.String())

	got := v1.Total()

	t.Log(got.String())
	t.Log(got.TotalWithRound(Boite))

	exp := amount.NewA([]quant.Q{
		{Inconnu, 21},
		{Pillule, 185},
	}...)

	ValEqualCheck(t, got, exp)
}
