package val

import (
	"github.com/olivier5741/stock-manager/item/unitval"
	"sort"
)

type Val struct {
	Vals map[string]unitval.T
}

func (v Val) String() (s string) {
	for _, u := range v.ValsByFactorDesc() {
		s += u.String() + ", "
	}
	return
}

func NewVal(units ...unitval.T) Val {
	vals := make(map[string]unitval.T, 0)
	for _, u := range units {
		vals[u.ID()] = u
	}
	return Val{vals}
}

func (v Val) Empty() (out Val) {
	vals := make(map[string]unitval.T, 0)
	for k, u := range v.Vals {
		vals[k] = u.NewSet(0)
	}
	return Val{vals}
}

func (v Val) Neg() Val {
	vals := make(map[string]unitval.T, 0)
	for _, val := range v.Values() {
		vals[val.ID()] = val.NewSet(-val.Val)
	}
	return Val{vals}
}

func Add(v1, v2 Val) (out Val) {

	out = v1.Copy()

	for _, v := range v2.Values() {
		if old, ok := out.Vals[v.ID()]; ok {
			v = unitval.Add(old, v)
		}
		out.Vals[v.ID()] = v
	}
	return
}

func Sub(v1, v2 Val) (out Val) {

	out = v1.Copy()
	out = stupidSubVal(out, v2.valsWithout())

	for _, v := range v2.ValsWithByFactorAsc() {
		if old, ok := out.Vals[v.ID()]; !ok {
			out.Vals[v.ID()] = unitval.T{v.Unit, -v.Val}
		} else {
			if old.Val >= v.Val || v.Val > out.TotalWith() {
				out.Vals[v.ID()] = unitval.Sub(old, v)
				continue
			}
			current := out.ValsWithByFactorAsc()
			needed := valsUntilAboveLimit(current, v.Total())
			sort.Sort(unitval.ByFactDesc(needed))
			missing := v.Total()
			for i, n := range needed {
				left := valsTotal(needed[i+1:])

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
				stupidSubVal(out, map[string]unitval.T{v.ID(): {v.Unit, missing}})
			}

		}
	}

	return
}

// Perhaps delete noWithout
func Diff(v1, v2 Val) (out Val, noWithout bool, diff int) {
	out = Sub(v1, v2).Redistribute()
	noWithout = len(out.valsWithout()) == 0
	diff = out.TotalWith()
	return
}

func (v Val) Redistribute() Val {
	out := v.valsWithout()
	var lasts []unitval.T
	left := 0
	for _, val := range v.ValsWithByFactorDesc() {
		out[val.ID()] = val.Empty()
		lasts = append(lasts, val.Empty())
		left += val.Total()
		for _, l := range lasts {
			out[l.ID()] = out[l.ID()].NewAdd(left / l.Fact)
			left = left % l.Fact
		}
	}
	return Val{out}
}

func (v Val) Total() (out Val) {
	out = Val{v.valsWithout()}
	with := v.ValsWithByFactorDesc()
	total := 0
	for _, val := range with {
		total += val.Total()
	}
	if len(with) > 0 {
		smallest := with[len(with)-1]
		out.Vals[smallest.ID()] =
			smallest.NewSet(total / smallest.Fact)
	}
	return
}

func (v Val) TotalWithRound(u unitval.Unit) (out unitval.TFloat) {
	total := 0.0
	for _, val := range v.valsWith() {
		total += float64(val.Total()) / float64(u.Fact)
	}
	return unitval.TFloat{u, total}
}

func (v Val) TotalWith() (total int) {
	for _, val := range v.valsWith() {
		total += val.Total()
	}
	return
}

func (v Val) Copy() Val {
	return NewVal(v.Values()...)
}

func (v Val) valsWithout() (out map[string]unitval.T) {
	_, out = v.ValuesFactFilter()
	return
}

func (v Val) valsWith() (out map[string]unitval.T) {
	out, _ = v.ValuesFactFilter()
	return
}

func (v Val) ValsWithByFactorAsc() []unitval.T {
	return Val{Vals: v.valsWith()}.ValsByFactorAsc()
}

func (v Val) ValsWithByFactorDesc() []unitval.T {
	return Val{Vals: v.valsWith()}.ValsByFactorDesc()
}

func stupidSubVal(v1 Val, vals map[string]unitval.T) (val Val) {
	val = v1.Copy()
	for _, v := range vals {
		if old, ok := val.Vals[v.ID()]; ok {
			val.Vals[v.ID()] = unitval.Sub(old, v)
		} else {
			val.Vals[v.ID()] = unitval.T{v.Unit, -v.Val}
		}
	}
	return
}

func valsUntilAboveLimit(vals []unitval.T, limit int) []unitval.T {
	var out []unitval.T
	for _, v := range vals {
		out = append(out, v)
		if valsTotal(out) >= limit {
			return out
		}
	}
	return out
}

func valsTotal(vals []unitval.T) (total int) {
	for _, v := range vals {
		total += v.Total()
	}
	return
}

func copyMap(originalMap map[string]unitval.T) map[string]unitval.T {
	newMap := make(map[string]unitval.T, 0)
	for k, v := range originalMap {
		newMap[k] = v
	}
	return newMap
}

func (v Val) ValuesFactFilter() (with map[string]unitval.T, without map[string]unitval.T) {
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

func (v Val) Values() []unitval.T {
	var list []unitval.T
	for _, val := range v.Vals {
		list = append(list, val)
	}
	return list
}

func (v Val) ValsByFactorDesc() []unitval.T {
	out := v.Values()
	sort.Sort(unitval.ByFactDesc(out))
	return out
}

func (v Val) ValsByFactorAsc() []unitval.T {
	out := v.Values()
	sort.Sort(unitval.ByFactAsc(out))
	return out
}
