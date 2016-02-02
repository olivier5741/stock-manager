package val

import (
	"github.com/olivier5741/stock-manager/item/unitval"
	"sort"
)

type T struct {
	Vals map[string]unitval.T
}

func (v T) String() string {
	var s string
	for _, u := range v.ValsByFactDesc() {
		s += u.String() + ", "
	}
	return s
}

func NewT(units ...unitval.T) T {
	vals := make(map[string]unitval.T, 0)
	for _, u := range units {
		vals[u.ID()] = u
	}
	return T{vals}
}

func (v T) Empty() T {
	vals := make(map[string]unitval.T, 0)
	for k, u := range v.Vals {
		vals[k] = u.NewSet(0)
	}
	return T{vals}
}

func (v T) Neg() T {
	vals := make(map[string]unitval.T, 0)
	for _, val := range v.Values() {
		vals[val.ID()] = val.NewSet(-val.Val)
	}
	return T{vals}
}

func Add(v1, v2 T) T {
	out := v1.Copy()
	for _, v := range v2.Values() {
		if old, ok := out.Vals[v.ID()]; ok {
			v = unitval.Add(old, v)
		}
		out.Vals[v.ID()] = v
	}
	return out
}

func Sub(v1, v2 T) T {
	out := v1.Copy()
	out = valByValSub(out, v2.valsWithout())
	for _, v := range v2.ValsWithByFactAsc() {
		if old, ok := out.Vals[v.ID()]; !ok {
			out.Vals[v.ID()] = unitval.T{v.Unit, -v.Val}
		} else {
			if old.Val >= v.Val || v.Val > out.TotalWith() {
				out.Vals[v.ID()] = unitval.Sub(old, v)
				continue
			}
			current := out.ValsWithByFactAsc()
			needed := unitval.SliceFromSliceTotalUntilLimit(current, v.Total())
			sort.Sort(unitval.ByFactDesc(needed))
			missing := v.Total()
			for i, n := range needed {
				left := unitval.SliceTotal(needed[i+1:])

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

				out.Vals[n.ID()] = n.NewSub(tosub)
				missing -= tosub * n.Fact
			}
			if missing > 0 {
				valByValSub(out, map[string]unitval.T{v.ID(): {v.Unit, missing}})
			}

		}
	}

	return out
}

// Perhaps delete noWithout
func Diff(v1, v2 T) (out T, noWithout bool, diff int) {
	out = Sub(v1, v2).Redistribute()
	noWithout = len(out.valsWithout()) == 0
	diff = out.TotalWith()
	return
}

func (v T) Redistribute() T {
	out := v.valsWithout()
	var lasts []unitval.T
	left := 0
	for _, val := range v.ValsWithByFactDesc() {
		out[val.ID()] = val.Empty()
		lasts = append(lasts, val.Empty())
		left += val.Total()
		for _, l := range lasts {
			out[l.ID()] = out[l.ID()].NewAdd(left / l.Fact)
			left = left % l.Fact
		}
	}
	return T{out}
}

func (v T) Total() T {
	out := T{v.valsWithout()}
	with := v.ValsWithByFactDesc()
	total := 0
	for _, val := range with {
		total += val.Total()
	}
	if len(with) > 0 {
		smallest := with[len(with)-1]
		out.Vals[smallest.ID()] =
			smallest.NewSet(total / smallest.Fact)
	}
	return out
}

func (v T) TotalWithRound(u unitval.Unit) unitval.TFloat {
	total := 0.0
	for _, val := range v.valsWith() {
		total += float64(val.Total()) / float64(u.Fact)
	}
	return unitval.TFloat{u, total}
}

func (v T) TotalWith() int {
	var t int
	for _, val := range v.valsWith() {
		t += val.Total()
	}
	return t
}

func (v T) Copy() T {
	return NewT(v.Values()...)
}

func (v T) valsWithout() map[string]unitval.T {
	_, out := v.ValsFactFilter()
	return out
}

func (v T) valsWith() map[string]unitval.T {
	out, _ := v.ValsFactFilter()
	return out
}

func (v T) ValsWithByFactAsc() []unitval.T {
	return T{Vals: v.valsWith()}.ValsByFactAsc()
}

func (v T) ValsWithByFactDesc() []unitval.T {
	return T{Vals: v.valsWith()}.ValsByFactDesc()
}

func valByValSub(v1 T, vals map[string]unitval.T) T {
	val := v1.Copy()
	for _, v := range vals {
		if old, ok := val.Vals[v.ID()]; ok {
			val.Vals[v.ID()] = unitval.Sub(old, v)
		} else {
			val.Vals[v.ID()] = unitval.T{v.Unit, -v.Val}
		}
	}
	return val
}

func (v T) ValsFactFilter() (with, without map[string]unitval.T) {
	with = make(map[string]unitval.T, 0)
	without = make(map[string]unitval.T, 0)
	for _, val := range v.Vals {
		if val.Fact != 0 {
			with[val.ID()] = val
		} else {
			without[val.ID()] = val
		}
	}
	return
}

func (v T) Values() []unitval.T {
	var list []unitval.T
	for _, val := range v.Vals {
		list = append(list, val)
	}
	return list
}

func (v T) ValsByFactDesc() []unitval.T {
	out := v.Values()
	sort.Sort(unitval.ByFactDesc(out))
	return out
}

func (v T) ValsByFactAsc() []unitval.T {
	out := v.Values()
	sort.Sort(unitval.ByFactAsc(out))
	return out
}
