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
func NewA(qs ...quant.Q) A {
	list := make(map[string]quant.Q, 0)
	for _, q := range qs {
		list[q.ID()] = q
	}
	return A{list}
}

// Empty creates a new amount based on am
// with all quantities value set to 0
func (am A) Empty() A {
	list := make(map[string]quant.Q, 0)
	for k, q := range am.quants {
		list[k] = q.NewSet(0)
	}
	return A{list}
}

// Neg creates a new amount based on am
// with all quantities value as the opposite
// of the
func (am A) Neg() A {
	list := make(map[string]quant.Q, 0)
	for _, q := range am.Quants() {
		list[q.ID()] = q.NewSet(-q.Val)
	}
	return A{list}
}

// Add creates a new amount by summing each quantities
// value of a1 and a2 and append the unmatched quantities
func Add(a1, a2 A) A {
	out := a1.Copy()
	for _, q := range a2.Quants() {
		if old, ok := out.quants[q.ID()]; ok {
			q = quant.Add(old, q)
		}
		out.quants[q.ID()] = q
	}
	return out
}

// Sub creates a new amount by subtracting a2 from a1
// TODO : more explanation
func Sub(a1, a2 A) A {
	out := a1.Copy()
	out = valByValSub(out, a2.quantsWithout())
	for _, q := range a2.QuantsWithByFactAsc() {
		if old, ok := out.quants[q.ID()]; !ok {
			out.quants[q.ID()] = quant.Q{q.Unit, -q.Val}
		} else {
			if old.Val >= q.Val || q.Val > out.TotalWith() {
				out.quants[q.ID()] = quant.Sub(old, q)
				continue
			}
			current := out.QuantsWithByFactAsc()
			needed := quant.TrimSliceOnTotal(current, q.Total())
			sort.Sort(quant.ByFactDesc(needed))
			missing := q.Total()
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
				valByValSub(out, map[string]quant.Q{q.ID(): {q.Unit, missing}})
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
	noWithout = len(out.quantsWithout()) == 0
	diff = out.TotalWith()
	return
}

// Redistribute creates an amount by redistributing
// parts of possible quantities values of am
// to a quantity with a higher factor
func (am A) Redistribute() A {
	list := am.quantsWithout()
	var lasts []quant.Q
	left := 0
	for _, q := range am.QuantsWithByFactDesc() {
		list[q.ID()] = q.Empty()
		lasts = append(lasts, q.Empty())
		left += q.Total()
		for _, l := range lasts {
			list[l.ID()] = list[l.ID()].NewAdd(left / l.Fact)
			left = left % l.Fact
		}
	}
	return A{list}
}

// Total creates an amount by redistributing all the quatities values
// to the quantity with the smallest factor (0 factor means no factor so
// it is not taken into account)
func (am A) Total() A {
	out := A{am.quantsWithout()}
	with := am.QuantsWithByFactDesc()
	t := 0
	for _, q := range with {
		t += q.Total()
	}
	if len(with) > 0 {
		smallest := with[len(with)-1]
		out.quants[smallest.ID()] =
			smallest.NewSet(t / smallest.Fact)
	}
	return out
}

// TotalWithRound returns TODO
func (am A) TotalWithRound(u quant.Unit) quant.QFloat {
	t := 0.0
	for _, q := range am.quantsWith() {
		t += float64(q.Total()) / float64(u.Fact)
	}
	return quant.QFloat{u, t}
}

// TotalWith returns TODO
func (am A) TotalWith() int {
	var t int
	for _, q := range am.quantsWith() {
		t += q.Total()
	}
	return t
}

// Copy creates a new amount based on am
func (am A) Copy() A {
	return NewA(am.Quants()...)
}

func (am A) quantsWithout() map[string]quant.Q {
	_, out := am.valsFactFilter()
	return out
}

func (am A) quantsWith() map[string]quant.Q {
	out, _ := am.valsFactFilter()
	return out
}

func (am A) QuantsWithByFactAsc() []quant.Q {
	return A{quants: am.quantsWith()}.QuantsByFactAsc()
}

func (am A) QuantsWithByFactDesc() []quant.Q {
	return A{quants: am.quantsWith()}.QuantsByFactDesc()
}

func valByValSub(am A, qs map[string]quant.Q) A {
	out := am.Copy()
	for _, q := range qs {
		if old, ok := out.quants[q.ID()]; ok {
			out.quants[q.ID()] = quant.Sub(old, q)
		} else {
			out.quants[q.ID()] = quant.Q{q.Unit, -q.Val}
		}
	}
	return out
}

func (am A) valsFactFilter() (with, without map[string]quant.Q) {
	with = make(map[string]quant.Q, 0)
	without = make(map[string]quant.Q, 0)
	for _, q := range am.quants {
		if q.Fact != 0 {
			with[q.ID()] = q
		} else {
			without[q.ID()] = q
		}
	}
	return
}

// Quants returns the quantities of a
func (am A) Quants() []quant.Q {
	var out []quant.Q
	for _, val := range am.quants {
		out = append(out, val)
	}
	return out
}

// QuantsMap return the quantities of a in a map
func (am A) QuantsMap() map[string]quant.Q {
	out := make(map[string]quant.Q)
	for _, val := range am.quants {
		out[val.Unit.ID()] = val
	}
	return out
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
