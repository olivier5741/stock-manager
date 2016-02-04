// Package amount provides operations  on an amount
// such as addition, subtraction ...
package amount

import (
	"github.com/olivier5741/stock-manager/item/quant"
	"sort"
)

// A represents an amount which consists of several quantities
// in the same system
type A struct {
	quants map[string]quant.Q
}

// String returns a comma seperated list of quantities
// sorted by descending factor
func (am A) String() string {
	var s string
	for _, u := range am.QuantsByFactDesc() {
		s += u.String() + ", "
	}
	return s
}

// NewA creates a new amount based on qs quantities
func NewA(units ...quant.Q) A {
	vals := make(map[string]quant.Q, 0)
	for _, u := range units {
		vals[u.ID()] = u
	}
	return A{vals}
}

// Empty creates a new amount based on am
// with all quantities value set to 0
func (am A) Empty() A {
	vals := make(map[string]quant.Q, 0)
	for k, u := range am.quants {
		vals[k] = u.NewSet(0)
	}
	return A{vals}
}

// Neg creates a new amount based on am
// with all quantities value as the opposite
// of the
func (am A) Neg() A {
	vals := make(map[string]quant.Q, 0)
	for _, val := range am.Quants() {
		vals[val.ID()] = val.NewSet(-val.Val)
	}
	return A{vals}
}

// Add creates a new amount by summing each quantities
// value of a1 and a2 and append the unmatched quantities
func Add(a1, a2 A) A {
	out := a1.Copy()
	for _, v := range a2.Quants() {
		if old, ok := out.quants[v.ID()]; ok {
			v = quant.Add(old, v)
		}
		out.quants[v.ID()] = v
	}
	return out
}

// Sub creates a new amount by subtracting a2 from a1
// TODO : more explanation
func Sub(a1, a2 A) A {
	out := a1.Copy()
	out = valByValSub(out, a2.valsWithout())
	for _, v := range a2.ValsWithByFactAsc() {
		if old, ok := out.quants[v.ID()]; !ok {
			out.quants[v.ID()] = quant.Q{v.Unit, -v.Val}
		} else {
			if old.Val >= v.Val || v.Val > out.TotalWith() {
				out.quants[v.ID()] = quant.Sub(old, v)
				continue
			}
			current := out.ValsWithByFactAsc()
			needed := quant.TrimSliceOnTotal(current, v.Total())
			sort.Sort(quant.ByFactDesc(needed))
			missing := v.Total()
			for i, n := range needed {
				left := quant.SliceTotal(needed[i+1:])

				if left == missing {
					continue
				}

				tosub := missing / n.Fact

				if left < missing {
					tosub = (missing - left) / n.Fact
					if tosub*n.Fact != missing-left {
						tosub++
					}
				}

				out.quants[n.ID()] = n.NewSub(tosub)
				missing -= tosub * n.Fact
			}
			if missing > 0 {
				valByValSub(out, map[string]quant.Q{v.ID(): {v.Unit, missing}})
			}

		}
	}

	return out
}

// Diff creates an amount resulting from the subtraction of a2 from a1
// and also returns if this amount as a single quantity and the value
// of this quantity
// TODO  Perhaps delete noWithout
func Diff(a1, a2 A) (out A, noWithout bool, diff int) {
	out = Sub(a1, a2).Redistribute()
	noWithout = len(out.valsWithout()) == 0
	diff = out.TotalWith()
	return
}

// Redistribute creates an amount by redistributing
// parts of possible quantities values of am
// to a quantity with a higher factor
func (am A) Redistribute() A {
	out := am.valsWithout()
	var lasts []quant.Q
	left := 0
	for _, val := range am.ValsWithByFactDesc() {
		out[val.ID()] = val.Empty()
		lasts = append(lasts, val.Empty())
		left += val.Total()
		for _, l := range lasts {
			out[l.ID()] = out[l.ID()].NewAdd(left / l.Fact)
			left = left % l.Fact
		}
	}
	return A{out}
}

// Total creates an amount by redistributing all the quatities values
// to the quantity with the smallest factor (0 factor means no factor so
// it is not taken into account)
func (am A) Total() A {
	out := A{am.valsWithout()}
	with := am.ValsWithByFactDesc()
	total := 0
	for _, val := range with {
		total += val.Total()
	}
	if len(with) > 0 {
		smallest := with[len(with)-1]
		out.quants[smallest.ID()] =
			smallest.NewSet(total / smallest.Fact)
	}
	return out
}

// TotalWithRound returns TODO
func (am A) TotalWithRound(u quant.Unit) quant.QFloat {
	total := 0.0
	for _, val := range am.valsWith() {
		total += float64(val.Total()) / float64(u.Fact)
	}
	return quant.QFloat{u, total}
}

// TotalWith returns TODO
func (am A) TotalWith() int {
	var t int
	for _, val := range am.valsWith() {
		t += val.Total()
	}
	return t
}

// Copy creates a new amount based on am
func (am A) Copy() A {
	return NewA(am.Quants()...)
}

func (am A) valsWithout() map[string]quant.Q {
	_, out := am.valsFactFilter()
	return out
}

func (am A) valsWith() map[string]quant.Q {
	out, _ := am.valsFactFilter()
	return out
}

func (am A) ValsWithByFactAsc() []quant.Q {
	return A{quants: am.valsWith()}.QuantsByFactAsc()
}

func (am A) ValsWithByFactDesc() []quant.Q {
	return A{quants: am.valsWith()}.QuantsByFactDesc()
}

func valByValSub(am A, vals map[string]quant.Q) A {
	val := am.Copy()
	for _, v := range vals {
		if old, ok := val.quants[v.ID()]; ok {
			val.quants[v.ID()] = quant.Sub(old, v)
		} else {
			val.quants[v.ID()] = quant.Q{v.Unit, -v.Val}
		}
	}
	return val
}

func (am A) valsFactFilter() (with, without map[string]quant.Q) {
	with = make(map[string]quant.Q, 0)
	without = make(map[string]quant.Q, 0)
	for _, val := range am.quants {
		if val.Fact != 0 {
			with[val.ID()] = val
		} else {
			without[val.ID()] = val
		}
	}
	return
}

// Quants returns the quantities of a
func (am A) Quants() []quant.Q {
	var list []quant.Q
	for _, val := range am.quants {
		list = append(list, val)
	}
	return list
}

// QuantsMap return the quantities of a in a map
func (am A) QuantsMap() map[string]quant.Q {
	list := make(map[string]quant.Q)
	for _, val := range am.quants {
		list[val.Unit.ID()] = val
	}
	return list
}

// QuantsByFactDesc returns the quantities of am by descending factor
func (am A) QuantsByFactDesc() []quant.Q {
	out := am.Quants()
	sort.Sort(quant.ByFactDesc(out))
	return out
}

// QuantsByFactAsc returns the quantities of am by ascending factor
func (am A) QuantsByFactAsc() []quant.Q {
	out := am.Quants()
	sort.Sort(quant.ByFactAsc(out))
	return out
}
