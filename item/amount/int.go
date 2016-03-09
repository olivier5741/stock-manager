// Package amount provides operations  on an amount
// such as addition, subtraction ...
package amount

import (
	"github.com/olivier5741/stock-manager/item/quant"
	"sort"
	"math"
	"math/big"
)

// Amount represents an amount which consists of several quantities
// in the same system
type Amount struct {
	quants map[string]quant.Quant
}

// String returns a comma seperated list of quantities
// sorted by descending factor
func (am Amount) String() string {
	var s string
	for _, u := range am.QuantsByFactDesc() {
		s += u.String() + ", "
	}
	return s
}

// NewA creates a new amount based on qs quantities
func NewAmount(qs ...quant.Quant) Amount {
	list := make(map[string]quant.Quant, 0)
	for _, q := range qs {
		list[q.ID()] = q
	}
	return Amount{list}
}

func FromStringSlice(l []string) Amount {
	// could put this in quant ...
	var quants []quant.Quant
	for i := 0; i < len(l); i = i + 2 {

		if l[i] == "" {
			continue // empty string, should be filtered before
		}

		//val,_ := new(big.Rat).SetString(l[i])
		quants = append(quants, quant.Quant{quant.NewUnit(l[i+1]), quant.StringToRat(l[i])})
	}
	return NewAmount(quants...)
}

// Empty creates a new amount based on am
// with all quantities value set to 0
func (am Amount) Empty() Amount {
	list := make(map[string]quant.Quant, 0)
	for k, q := range am.quants {
		list[k] = q.NewSet(&big.Rat{})
	}
	return Amount{list}
}

// Neg creates a new amount based on am
// with all quantities value as the opposite
// of the
func (am Amount) Neg() Amount {
	list := make(map[string]quant.Quant, 0)
	for _, q := range am.Quants() {
		list[q.ID()] = q.NewSet(new(big.Rat).Neg(q.Val))
	}
	return Amount{list}
}

// Add creates a new amount by summing each quantities
// value of a1 and a2 and append the unmatched quantities
func Add(a1, a2 Amount) Amount {
	out := a1.Copy()
	for _, q := range a2.Quants() {
		if old, ok := out.quants[q.ID()]; ok {
			q = quant.Add(old, q)
		}
		out.quants[q.ID()] = q
	}
	return out
}

func IntDiv(i1, i2 *big.Rat) *big.Rat {
	f,_ := new(big.Rat).Quo(i1,i2).Float64()
	return new(big.Rat).SetFloat64(math.Floor(f))
}

// Sub creates a new amount by subtracting a2 from a1
// TODO : WHEN QUANTITIES NOT PRESENT AFTER, NOT ACTING VERY WELL
// Il faut rajouter les quantité inconnues de a2 dans a1
func Sub(a1, a2 Amount) Amount {
	out := a1.Copy()
	out = valByValSub(out, a2.quantsWithout())
	for _, q := range a2.QuantsWithByFactAsc() {

		old, ok := out.quants[q.ID()]

		if !ok {
			out.quants[q.ID()] = quant.Quant{q.Unit, &big.Rat{}}
			continue
		}

		if old.Val.Cmp(q.Val) >= 0 || q.Val.Cmp(out.TotalWith()) > 0 {
			out.quants[q.ID()] = quant.Sub(old, q)
			continue
		}
		current := out.QuantsWithByFactAsc()
		needed := quant.TrimSliceOnTotal(current, q.Total())
		sort.Sort(quant.ByFactDesc(needed))
		missing := q.Total()
		for i, n := range needed {
			
			left := quant.SliceTotal(needed[i+1:])

			if left.Cmp(missing) == 0 {
				continue
			}

			tosub := IntDiv(missing,n.Fact)

			if left.Cmp(missing) < 0 {
				tosub = IntDiv(new(big.Rat).Sub(missing,left),n.Fact)
				if new(big.Rat).Mul(tosub,n.Fact).Cmp(new(big.Rat).Sub(missing,left)) != 0 {
					tosub = new(big.Rat).Add(tosub,big.NewRat(1,1))
				}
			}

			out.quants[n.ID()] = n.NewSub(tosub)
			missing.Sub(missing,new(big.Rat).Mul(tosub,n.Fact))
		}
		if missing.Cmp(&big.Rat{}) > 0 {
			valByValSub(out, map[string]quant.Quant{q.ID(): {q.Unit, missing}})
		}		
	}

	return out
}

// Diff creates an amount resulting from the subtraction of a2 from a1
// and also returns if this amount as a single quantity and the value
// of this quantity
// TODO  Perhaps delete noWithout
func Diff(a1, a2 Amount) (out Amount, noWithout bool, diff *big.Rat) {
	out = Sub(a1, a2).Redistribute()
	noWithout = len(out.quantsWithout()) == 0
	diff = out.TotalWith()
	return
}

// Redistribute creates an amount by redistributing
// parts of possible quantities values of am
// to a quantity with a higher factor
// TODO : Rational and modulo not playing nice I think
func (am Amount) Redistribute() Amount {
	list := am.quantsWithout()
	var lasts []quant.Quant
	left := &big.Rat{}
	for _, q := range am.QuantsWithByFactDesc() {
		list[q.ID()] = q.Empty()
		lasts = append(lasts, q.Empty())
		left.Add(left,q.Total())
		for _, l := range lasts {
			a,_ := new(big.Rat).Quo(left,l.Fact).Float64()
			div := new(big.Rat).SetFloat64(math.Floor(a))
			list[l.ID()] = list[l.ID()].NewAdd(div)
			left = new(big.Rat).Sub(left,new(big.Rat).Mul(div,l.Fact))
		}
	}
	return Amount{list}
}

// Total creates an amount by redistributing all the quatities values
// to the quantity with the smallest factor (0 factor means no factor so
// it is not taken into account)
func (am Amount) Total() Amount {
	out := Amount{am.quantsWithout()}
	with := am.QuantsWithByFactDesc()
	t := &big.Rat{}
	for _, q := range with {
		t.Add(t,q.Total())
	}
	if len(with) > 0 {
		smallest := with[len(with)-1]
		out.quants[smallest.ID()] =
			smallest.NewSet(new(big.Rat).Quo(t,smallest.Fact))
	}
	return out
}

// TotalWithRound returns TODO
func (am Amount) TotalWithRound1(u quant.Unit) *big.Rat { // rename
	t := &big.Rat{}
	for _, q := range am.quantsWith() {
		t = new(big.Rat).Add(t,new(big.Rat).Quo(q.Total(),u.Fact))
	}
	f,_ := t.Float64()
	return  new(big.Rat).SetFloat64(math.Ceil(f))
}

func (am Amount) TotalWithRound(u quant.Unit) Amount{
	return NewAmount(quant.Quant{u,am.TotalWithRound1(u)})
}

// TotalWith returns TODO
func (am Amount) TotalWith() *big.Rat {
	t := &big.Rat{}
	for _, q := range am.quantsWith() {
		t = new(big.Rat).Add(t,q.Total())
	}
	return t
}

// Copy creates a new amount based on am
func (am Amount) Copy() Amount {
	return NewAmount(am.Quants()...)
}

func (am Amount) quantsWithout() map[string]quant.Quant {
	_, out := am.valsFactFilter()
	return out
}

func (am Amount) quantsWith() map[string]quant.Quant {
	out, _ := am.valsFactFilter()
	return out
}

func (am Amount) QuantsWithByFactAsc() []quant.Quant {
	return Amount{quants: am.quantsWith()}.QuantsByFactAsc()
}

func (am Amount) QuantsWithByFactDesc() []quant.Quant {
	return Amount{quants: am.quantsWith()}.QuantsByFactDesc()
}

func valByValSub(am Amount, qs map[string]quant.Quant) Amount {
	out := am.Copy()
	for _, q := range qs {
		if old, ok := out.quants[q.ID()]; ok {
			out.quants[q.ID()] = quant.Sub(old, q)
		} else {
			out.quants[q.ID()] = quant.Quant{q.Unit, new(big.Rat).Neg(q.Val)}
		}
	}
	return out
}

func (am Amount) valsFactFilter() (with, without map[string]quant.Quant) {
	with = make(map[string]quant.Quant, 0)
	without = make(map[string]quant.Quant, 0)
	for _, q := range am.quants {
		if q.Fact.Cmp(&big.Rat{}) != 0 {
			with[q.ID()] = q
		} else {
			without[q.ID()] = q
		}
	}
	return
}

// Quants returns the quantities of a
func (am Amount) Quants() []quant.Quant {
	var out []quant.Quant
	for _, val := range am.quants {
		out = append(out, val)
	}
	return out
}

// QuantsMap return the quantities of a in a map
func (am Amount) QuantsMap() map[string]quant.Quant {
	out := make(map[string]quant.Quant)
	for _, val := range am.quants {
		out[val.Unit.ID()] = val
	}
	return out
}

// QuantsByFactDesc returns the quantities of am by descending factor
func (am Amount) QuantsByFactDesc() []quant.Quant {
	out := am.Quants()
	sort.Sort(quant.ByFactDesc(out))
	return out
}

// QuantsByFactAsc returns the quantities of am by ascending factor
func (am Amount) QuantsByFactAsc() []quant.Quant {
	out := am.Quants()
	sort.Sort(quant.ByFactAsc(out))
	return out
}
