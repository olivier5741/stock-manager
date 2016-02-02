package unitval

import (
	"strconv"
)

type T struct {
	Unit
	Val int
}

func (u T) String() string {
	return strconv.Itoa(u.Val) + " " + u.Unit.String()
}

func (u T) Total() int {
	return u.Fact * u.Val
}

func (u T) Empty() T {
	return T{u.Unit, 0}
}

func (u T) NewAdd(add int) T {
	return T{u.Unit, u.Val + add}
}

func (u T) NewSub(sub int) T {
	return T{u.Unit, u.Val - sub}
}

func (u T) NewSet(set int) T {
	return T{u.Unit, set}
}

func Add(u1, u2 T) T {
	return T{u1.Unit, u1.Val + u2.Val}
}

func Sub(u1, u2 T) T {
	return T{u1.Unit, u1.Val - u2.Val}
}

func MapToSlice(m map[string]T) []T {
	var out []T
	for _, v := range m {
		out = append(out, v)
	}
	return out
}

func SliceTotal(vals []T) int {
	var total int
	for _, v := range vals {
		total += v.Total()
	}
	return total
}

func SliceFromSliceTotalUntilLimit(vals []T, lim int) []T {
	var out []T
	for _, v := range vals {
		out = append(out, v)
		if SliceTotal(out) >= lim {
			return out
		}
	}
	return out
}

func CopyMap(m map[string]T) map[string]T {
	n := make(map[string]T, 0)
	for k, v := range m {
		n[k] = v
	}
	return n
}

type ByFactDesc []T

func (f ByFactDesc) Len() int           { return len(f) }
func (f ByFactDesc) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }
func (f ByFactDesc) Less(i, j int) bool { return f[i].Fact > f[j].Fact }

type ByFactAsc []T

func (f ByFactAsc) Len() int           { return len(f) }
func (f ByFactAsc) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }
func (f ByFactAsc) Less(i, j int) bool { return f[i].Fact < f[j].Fact }
