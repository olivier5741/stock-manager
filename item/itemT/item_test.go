package itemT

import (
	"github.com/olivier5741/stock-manager/item/amount"
	"github.com/olivier5741/stock-manager/item/quant"
	"testing"
	"math/big"
)

func TestAddVal(t *testing.T) {
	v1 := amount.NewAmount([]quant.Quant{
		{Pillule, new(big.Rat).SetInt64(12)},
		{Tablette, new(big.Rat).SetInt64(2)},
		{Boite, new(big.Rat).SetInt64(2)},
		{Carton, new(big.Rat).SetInt64(1)},
	}...)

	t.Log(v1.String())

	v2 := amount.NewAmount([]quant.Quant{
		{Tablette, new(big.Rat).SetInt64(3)},
		{Boite, new(big.Rat).SetInt64(2)},
	}...)

	t.Log(v2.String())

	got := amount.Add(v1, v2)

	t.Log(got.String())
	t.Log(got.TotalWithRound(Carton))

	exp := amount.NewAmount([]quant.Quant{
		{Pillule, new(big.Rat).SetInt64(12)},
		{Tablette, new(big.Rat).SetInt64(5)},
		{Boite, new(big.Rat).SetInt64(4)},
		{Carton, new(big.Rat).SetInt64(1)},
	}...)

	ValEqualCheck(t, got, exp)
}

func TestSub(t *testing.T) {
	v1 := amount.NewAmount([]quant.Quant{
		{Pillule, new(big.Rat).SetInt64(12)},
		{Tablette, new(big.Rat).SetInt64(2)},
		{Boite, new(big.Rat).SetInt64(2)},
		{Carton, new(big.Rat).SetInt64(1)},
	}...)

	t.Log(v1.String())

	v2 := amount.NewAmount([]quant.Quant{
		{Pillule, new(big.Rat).SetInt64(16)},
		{Tablette, new(big.Rat).SetInt64(3)},
		{Boite, new(big.Rat).SetInt64(2)},
	}...)

	t.Log(v2.String())

	got := amount.Sub(v1, v2)

	t.Log(got.String())
	t.Log(got.TotalWithRound(Carton))

	exp := amount.NewAmount([]quant.Quant{
		{Pillule, new(big.Rat).SetInt64(11)},
		{Tablette, new(big.Rat).SetInt64(1)},
		{Boite, new(big.Rat).SetInt64(9)},
		{Carton, new(big.Rat).SetInt64(0)},
	}...)

	ValEqualCheck(t, got, exp)
}

func TestDiff(t *testing.T) {
	v1 := amount.NewAmount([]quant.Quant{
		{Pillule, new(big.Rat).SetInt64(12)},
		{Tablette, new(big.Rat).SetInt64(40)},
		{Boite, new(big.Rat).SetInt64(2)},
		{Carton, new(big.Rat).SetInt64(1)},
	}...)

	t.Log(v1.String())

	v2 := amount.NewAmount([]quant.Quant{
		{Pillule, new(big.Rat).SetInt64(16)},
		{Tablette, new(big.Rat).SetInt64(3)},
		{Boite, new(big.Rat).SetInt64(2)},
	}...)

	t.Log(v2.String())

	got, _, _ := amount.Diff(v1, v2)

	t.Log(got.String())
	t.Log(got.TotalWithRound(Carton))

	exp := amount.NewAmount([]quant.Quant{
		{Carton, new(big.Rat).SetInt64(2)},
		{Boite, new(big.Rat).SetInt64(2)},
		{Tablette, new(big.Rat).SetInt64(0)},
		{Pillule, new(big.Rat).SetInt64(11)},
	}...)

	ValEqualCheck(t, got, exp)
}

func TestRedistribute(t *testing.T) {
	// l'unité principale ne peut valoir 0

	v1 := amount.NewAmount([]quant.Quant{
		{Pillule, new(big.Rat).SetInt64(50)},
		{Tablette, new(big.Rat).SetInt64(3)},
		{Boite, new(big.Rat).SetInt64(2)},
	}...)

	t.Log(v1.String())

	got := v1.Redistribute()

	t.Log(got.String())
	t.Log(got.TotalWithRound(Boite))

	exp := amount.NewAmount([]quant.Quant{
		{Pillule, new(big.Rat).SetInt64(5)},
		{Tablette, new(big.Rat).SetInt64(0)},
		{Boite, new(big.Rat).SetInt64(4)},
	}...)

	ValEqualCheck(t, got, exp)
}

func TestTotal(t *testing.T) {
	v1 := amount.NewAmount([]quant.Quant{
		{Inconnu, new(big.Rat).SetInt64(21)},
		{Pillule, new(big.Rat).SetInt64(50)},
		{Tablette, new(big.Rat).SetInt64(3)},
		{Boite, new(big.Rat).SetInt64(2)},
	}...)

	t.Log(v1.String())

	got := v1.Total()

	t.Log(got.String())
	t.Log(got.TotalWithRound(Boite))

	exp := amount.NewAmount([]quant.Quant{
		{Inconnu, new(big.Rat).SetInt64(21)},
		{Pillule, new(big.Rat).SetInt64(185)},
	}...)

	ValEqualCheck(t, got, exp)
}
