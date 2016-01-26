package item

import (
	"fmt"
	//	"log"
	"sort"
	"strconv"
)

type Ider interface {
	Id() string
}

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
	return u.Name
}

type Val struct {
	Main string // Main fact should always be common factor of all the other unit.Fact
	Vals map[string]UnitVal
}

func (v Val) String() string {
	// TO DO : print all value
	return strconv.Itoa(int(v.Vals[v.Main].Val))
}

type Prod string

func (p Prod) String() string {
	return string(p)
}

func NewVal(units ...string) Val {
	vals := make(map[string]UnitVal, 0)
	for _, u := range units {
		vals[u] = UnitVal{Unit{u, 0}, 0}
	}
	return Val{units[0], vals}
}

func NewValFromUnits(units ...Unit) Val {
	vals := make(map[string]UnitVal, 0)
	for _, u := range units {
		vals[u.Id()] = UnitVal{u, 0}
	}
	return Val{units[0].Name, vals}
}

func NewValFromUnitVals(units ...UnitVal) Val {
	vals := make(map[string]UnitVal, 0)
	for _, u := range units {
		vals[u.Id()] = u
	}
	return Val{units[0].Name, vals}
}

func AddVal(v1, v2 Val) (out Val, err error) {
	if v1.Main != v2.Main {
		err = fmt.Errorf("Cannot add two Val with different Main")
	}

	out = CopyVal(v1)

	for _, v := range v2.Values() {
		if old, ok := out.Vals[v.Id()]; ok {
			v, err = AddUnitVal(old, v)
		}
		out.Vals[v.Id()] = v
	}
	return
}

func SubVal(v1, v2 Val) (val Val, err error) {
	if v1.Main != v2.Main {
		err = fmt.Errorf("Cannot sub two Val with different Main")
	}

	val = CopyVal(v1)
	with, without := v2.ValuesFactFilter()
	val, err = stupidSubVal(val, without)

	asc := mapToSlice(with)
	sort.Sort(ByFactorAsc(asc))

	for _, v := range asc {
		if old, ok := val.Vals[v.Id()]; ok {
			if old.Val >= v.Val || v.Val > val.TotalWith() {
				val.Vals[v.Id()], err = SubUnitVal(old, v)
			} else {
				with1, _ := val.ValuesFactFilter()
				current := mapToSlice(with1)
				sort.Sort(ByFactorAsc(current))
				needed := valsUntilAboveLimit(current, v.Val*v.Fact)
				sort.Sort(ByFactorDesc(needed))
				missing := v.Val * v.Fact
				for i, n := range needed {
					left := valsTotal(needed[i+1:])
					if left < missing {
						tosub := ((missing - left) / n.Fact)
						if (tosub * n.Fact) != (missing - left) {
							tosub += 1
						}
						val.Vals[n.Id()], _ = SubUnitVal(n, UnitVal{n.Unit, tosub})
						missing = missing - (tosub * n.Fact)
					} else if left > missing {
						tosub := missing / n.Fact
						val.Vals[n.Id()], _ = AddUnitVal(n, UnitVal{n.Unit, -tosub})
						missing -= tosub * n.Fact
					}
				}
				if missing > 0 {
					stupidSubVal(val, map[string]UnitVal{v.Id(): UnitVal{v.Unit, missing}})
				}
			}
		} else {
			val.Vals[v.Id()] = UnitVal{v.Unit, -v.Val}
		}
	}

	return
}

func (v Val) valsWithout() (out map[string]UnitVal) {
	_, out = v.ValuesFactFilter()
	return
}

func (v Val) valsWith() (out map[string]UnitVal) {
	out, _ = v.ValuesFactFilter()
	return
}

func (v Val) valsWithByFactorDesc() []UnitVal {
	return Val{Vals: v.valsWith()}.ValsByFactorDesc()
}

func (v Val) Redistribute() Val {
	out := v.valsWithout()
	lasts := make([]UnitVal, 0)
	left := 0
	for _, inVal := range v.valsWith() {
		out[inVal.Id()] = NewUnitValSet(inVal, 0)
		lasts = append(lasts, NewUnitValSet(inVal, 0))
		left += inVal.Total()
		for _, l := range lasts {
			out[l.Id()] = NewUnitValAdd(out[l.Id()], left/l.Fact)
			left = left % l.Fact
		}
	}
	return Val{v.Main, out}
}

func (v Val) Total() Val {
	// total no using remainder
	with, without := v.ValuesFactFilter()
	in := Val{Vals: with}.ValsByFactorDesc() // Refactor this
	out := Val{v.Main, without}
	total := 0
	for _, val := range in {
		total += val.Val * val.Fact
	}
	if len(in) > 0 {
		smallest := in[len(in)-1]
		out.Vals[smallest.Id()] = UnitVal{smallest.Unit, total / smallest.Fact}
	}
	return out
}

func (v Val) TotalWith() int {
	total := 0
	with, _ := v.ValuesFactFilter()
	for _, val := range with {
		total += val.Fact * val.Val
	}
	return total
}

func CopyVal(v Val) Val {
	out := NewValFromUnitVals(v.Values()...)
	out.Main = v.Main
	return out
}

func stupidSubVal(v1 Val, vals map[string]UnitVal) (val Val, err error) {
	val = CopyVal(v1)
	for _, v := range vals {
		if old, ok := val.Vals[v.Id()]; ok {
			val.Vals[v.Id()], err = SubUnitVal(old, v)
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

func valsTotal(vals []UnitVal) int {
	// should all have a fact != 0
	total := 0
	for _, v := range vals {
		total += v.Fact * v.Val
	}
	return total
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

func (u UnitVal) Total() int {
	return u.Fact * u.Val
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

func AddUnitVal(u1, u2 UnitVal) (out UnitVal, err error) {
	out, err = NewUnitValMergeInit(u1, u2)
	out.Val = u1.Val + u2.Val
	return
}

func SubUnitVal(u1, u2 UnitVal) (out UnitVal, err error) {
	out, err = NewUnitValMergeInit(u1, u2)
	out.Val = u1.Val - u2.Val
	return
}

func NewUnitValMergeInit(u1, u2 UnitVal) (out UnitVal, err error) {
	if u1.Unit.Name != u2.Unit.Name {
		err = fmt.Errorf("Cannot merge two UnitVal with different Name")
	}

	unit := Unit{Name: u1.Name}

	if u1.Unit.Fact == 0 {
		unit.Fact = u2.Unit.Fact
	} else if u2.Unit.Fact != 0 && u1.Unit.Fact != u2.Unit.Fact {
		err = fmt.Errorf("Cannot merge two UnitVal with different Fact")
	} else {
		unit.Fact = u1.Unit.Fact
	}

	out = UnitVal{Unit: unit}

	return
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
	sort.Sort(ByFactorAsc(out))
	return out
}

// func (its Items) Add(adds Items) {
// 	for key, add := range adds {
// 		if it, ok := its[key]; ok {
// 			add.Val.T = it.Val.T + add.Val.T
// 		}
// 		its[key] = add
// 	}
// }

// func (its Items) Sub(subs Items) {
// 	for key, sub := range subs {
// 		if it, ok := its[key]; ok {
// 			sub.Val.T = it.Val.T - sub.Val.T
// 		} else {
// 			sub.Val.T = -sub.Val.T
// 		}
// 		its[key] = sub
// 	}
// }

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

// // I don't think I'm using this function
// func (origs Items) Copy() (out Items) {
// 	out = make(Items, len(origs))
// 	for key, orig := range origs {
// 		out[key] = orig
// 	}
// 	return
// }
