package quant

import (
	"strconv"
)

// Q represents a quantity with unit and value (int)
type Q struct {
	Unit
	Val int
}

// String returns the value and unit of the quantity
func (q Q) String() string {
	return strconv.Itoa(q.Val) + " " + q.Unit.String()
}

// Total returns the unit factor by the value of the quantity
func (q Q) Total() int {
	return q.Fact * q.Val
}

// Empty creates a quantity based on q
func (q Q) Empty() Q {
	return Q{q.Unit, 0}
}

// NewAdd creates a quantity based on q
// and setsg its value to the sum of q value and add
func (q Q) NewAdd(add int) Q {
	return Q{q.Unit, q.Val + add}
}

// NewSub creates a quantity based on q
// and sets its value to the subtraction of sub to q value
func (q Q) NewSub(sub int) Q {
	return Q{q.Unit, q.Val - sub}
}

// NewSet creates a quantity based on q
// and sets its value to set
func (q Q) NewSet(set int) Q {
	return Q{q.Unit, set}
}

// Add creates a quantity based on q1
// and sets its value to the sum of q1 value and q2 value
func Add(q1, q2 Q) Q {
	return Q{q1.Unit, q1.Val + q2.Val}
}

// Sub creates a quantity based on q2
// and sets its value to the subtraction of q2 value to q1 value
func Sub(q1, q2 Q) Q {
	return Q{q1.Unit, q1.Val - q2.Val}
}

// MapToSlice returns a slice quantities from a map of  quantities
func MapToSlice(m map[string]Q) []Q {
	var out []Q
	for _, v := range m {
		out = append(out, v)
	}
	return out
}

// SliceTotal returns the sum of all quantities total
func SliceTotal(qs []Q) int {
	var total int
	for _, v := range qs {
		total += v.Total()
	}
	return total
}

// TrimSliceOnTotal trims quantities until the sum of their total reach lim
func TrimSliceOnTotal(qs []Q, lim int) []Q {
	var out []Q
	for _, v := range qs {
		out = append(out, v)
		if SliceTotal(out) >= lim {
			return out
		}
	}
	return out
}

// CopyMap copies a map of quantities
func CopyMap(m map[string]Q) map[string]Q {
	n := make(map[string]Q, 0)
	for k, v := range m {
		n[k] = v
	}
	return n
}

// ByFactDesc sorts quantities by descending unit factor
type ByFactDesc []Q

func (f ByFactDesc) Len() int           { return len(f) }
func (f ByFactDesc) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }
func (f ByFactDesc) Less(i, j int) bool { return f[i].Fact > f[j].Fact }

// ByFactAsc sorts quantities by ascending unit factor
type ByFactAsc []Q

func (f ByFactAsc) Len() int           { return len(f) }
func (f ByFactAsc) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }
func (f ByFactAsc) Less(i, j int) bool { return f[i].Fact < f[j].Fact }
