package item

import (
	"sort"
	"strconv"
)

type Items map[string]Item

type Item struct {
	Prod Prod
	Val  Val
}

func (i Item) String() string {
	return i.Prod.String() + ": " + i.Val.String()
}

type UnitVal struct {
	Unit
	Val int
}

type Unit struct {
	Name string
	Fact int
}

func (u Unit) Id() string {
	return u.String()
}

func (u Unit) String() string {
	return u.Name + "(" + strconv.Itoa(u.Fact) + ")"
}

type Val struct {
	Vals map[string]UnitVal
}

func (v Val) String() (s string) {
	for _, u := range v.ValsByFactorDesc() {
		s += u.String() + ", "
	}
	return
}

type Prod string

func (p Prod) String() string {
	return string(p)
}

func NewVal(units ...UnitVal) Val {
	vals := make(map[string]UnitVal, 0)
	for _, u := range units {
		vals[u.Id()] = u
	}
	return Val{vals}
}

func NegVal(v Val) Val {
	vals := make(map[string]UnitVal, 0)
	for _, val := range v.Values() {
		vals[val.Id()] = NewUnitValSet(val, -val.Val)
	}
	return Val{vals}
}

func AddVal(v1, v2 Val) (out Val) {

	out = CopyVal(v1)

	for _, v := range v2.Values() {
		if old, ok := out.Vals[v.Id()]; ok {
			v = AddUnitVal(old, v)
		}
		out.Vals[v.Id()] = v
	}
	return
}

func SubVal(v1, v2 Val) (out Val) {

	out = CopyVal(v1)

	out = CopyVal(v1)
	out = stupidSubVal(out, v2.valsWithout())

	for _, v := range v2.valsWithByFactorAsc() {
		if old, ok := out.Vals[v.Id()]; !ok {
			out.Vals[v.Id()] = UnitVal{v.Unit, -v.Val}
		} else {
			if old.Val >= v.Val || v.Val > out.TotalWith() {
				out.Vals[v.Id()] = SubUnitVal(old, v)
				continue
			}
			current := out.valsWithByFactorAsc()
			needed := valsUntilAboveLimit(current, v.Total())
			sort.Sort(ByFactorDesc(needed))
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
						tosub += 1
					}
				}

				out.Vals[n.Id()] = NewUnitValSub(n, tosub)
				missing -= tosub * n.Fact
			}
			if missing > 0 {
				stupidSubVal(out, map[string]UnitVal{v.Id(): {v.Unit, missing}})
			}

		}
	}

	return
}

func (v Val) Redistribute() Val {
	out := v.valsWithout()
	lasts := make([]UnitVal, 0)
	left := 0
	for _, val := range v.valsWithByFactorDesc() {
		out[val.Id()] = NewUnitValInit(val)
		lasts = append(lasts, NewUnitValInit(val))
		left += val.Total()
		for _, l := range lasts {
			out[l.Id()] = NewUnitValAdd(out[l.Id()], left/l.Fact)
			left = left % l.Fact
		}
	}
	return Val{out}
}

func (v Val) Total() (out Val) {
	out = Val{v.valsWithout()}
	with := v.valsWithByFactorDesc()
	total := 0
	for _, val := range with {
		total += val.Total()
	}
	if len(with) > 0 {
		smallest := with[len(with)-1]
		out.Vals[smallest.Id()] =
			NewUnitValSet(smallest, total/smallest.Fact)
	}
	return
}

func (v Val) TotalWith() (total int) {
	for _, val := range v.valsWith() {
		total += val.Total()
	}
	return
}

func CopyVal(v Val) Val {
	return NewVal(v.Values()...)
}

func (v Val) valsWithout() (out map[string]UnitVal) {
	_, out = v.ValuesFactFilter()
	return
}

func (v Val) valsWith() (out map[string]UnitVal) {
	out, _ = v.ValuesFactFilter()
	return
}

func (v Val) valsWithByFactorAsc() []UnitVal {
	return Val{Vals: v.valsWith()}.ValsByFactorAsc()
}

func (v Val) valsWithByFactorDesc() []UnitVal {
	return Val{Vals: v.valsWith()}.ValsByFactorDesc()
}

func stupidSubVal(v1 Val, vals map[string]UnitVal) (val Val) {
	val = CopyVal(v1)
	for _, v := range vals {
		if old, ok := val.Vals[v.Id()]; ok {
			val.Vals[v.Id()] = SubUnitVal(old, v)
		} else {
			val.Vals[v.Id()] = UnitVal{v.Unit, -v.Val}
		}
	}
	return
}

func valsUntilAboveLimit(vals []UnitVal, limit int) []UnitVal {
	out := make([]UnitVal, 0)
	for _, v := range vals {
		out = append(out, v)
		if valsTotal(out) >= limit {
			return out
		}
	}
	return out
}

func valsTotal(vals []UnitVal) (total int) {
	for _, v := range vals {
		total += v.Total()
	}
	return
}

func copyMap(originalMap map[string]UnitVal) map[string]UnitVal {
	newMap := make(map[string]UnitVal, 0)
	for k, v := range originalMap {
		newMap[k] = v
	}
	return newMap
}

func (v Val) ValuesFactFilter() (with map[string]UnitVal, without map[string]UnitVal) {
	with = make(map[string]UnitVal, 0)
	without = make(map[string]UnitVal, 0)
	for _, val := range v.Vals {
		if val.Fact != 0 {
			with[val.Id()] = val
		} else {
			without[val.Id()] = val
		}
	}
	return
}

func (v Val) Values() []UnitVal {
	list := make([]UnitVal, 0)
	for _, val := range v.Vals {
		list = append(list, val)
	}
	return list
}

func (v Val) ValsByFactorDesc() []UnitVal {
	out := v.Values()
	sort.Sort(ByFactorDesc(out))
	return out
}

func (v Val) ValsByFactorAsc() []UnitVal {
	out := v.Values()
	sort.Sort(ByFactorAsc(out))
	return out
}

func (u UnitVal) String() string {
	return strconv.Itoa(u.Val) + " " + u.Unit.String()
}

func (u UnitVal) Total() int {
	return u.Fact * u.Val
}

func NewUnitValInit(prev UnitVal) UnitVal {
	return UnitVal{prev.Unit, 0}
}

func NewUnitValAdd(prev UnitVal, add int) UnitVal {
	return UnitVal{prev.Unit, prev.Val + add}
}

func NewUnitValSub(prev UnitVal, sub int) UnitVal {
	return UnitVal{prev.Unit, prev.Val - sub}
}

func NewUnitValSet(prev UnitVal, set int) UnitVal {
	return UnitVal{prev.Unit, set}
}

func AddUnitVal(u1, u2 UnitVal) UnitVal {
	return UnitVal{u1.Unit, u1.Val + u2.Val}
}

func SubUnitVal(u1, u2 UnitVal) UnitVal {
	return UnitVal{u1.Unit, u1.Val - u2.Val}
}

type ByFactorDesc []UnitVal

func (f ByFactorDesc) Len() int           { return len(f) }
func (f ByFactorDesc) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }
func (f ByFactorDesc) Less(i, j int) bool { return f[i].Fact > f[j].Fact }

type ByFactorAsc []UnitVal

func (f ByFactorAsc) Len() int           { return len(f) }
func (f ByFactorAsc) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }
func (f ByFactorAsc) Less(i, j int) bool { return f[i].Fact < f[j].Fact }

func mapToSlice(m map[string]UnitVal) []UnitVal {
	out := make([]UnitVal, 0)
	for _, v := range m {
		out = append(out, v)
	}
	return out
}

func (its Items) Add(adds Items) {
	for key, add := range adds {
		if it, ok := its[key]; ok {
			add.Val = AddVal(it.Val, add.Val)
		}
		its[key] = add
	}
}

func (its Items) Sub(subs Items) {
	for key, sub := range subs {
		if it, ok := its[key]; ok {
			sub.Val = SubVal(it.Val, sub.Val)
		} else {
			sub.Val = NegVal(sub.Val)
		}
		its[key] = sub
	}
}

// func (its Items) Missing(exps Items) (out Items) {
// 	out = map[string]Item{}
// 	for key, exp := range exps {
// 		if it, ok := its[key]; ok {
// 			if diff := exp.Val.T - it.Val.T; diff > 0 {
// 				out[key] = Item{it.Prod, Val{diff}}
// 			}
// 		} else {
// 			out[key] = Item{exp.Prod, Val{exp.Val.T}}
// 		}
// 	}
// 	return
// }

// I don't think I'm using this function
func (origs Items) Copy() (out Items) {
	out = make(Items, len(origs))
	for key, orig := range origs {
		out[key] = orig
	}
	return
}
