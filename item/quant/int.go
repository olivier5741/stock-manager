package quant

import (
	"strconv"
	"strings"
)

// Unit a unit with a name and a factor (must be a factor of smallest
// unit factor in the system)
//
// A system is for instance all the units a particular product can have
// ibuprofen : pill, tab, box, ...
type Unit struct {
	Name string
	Fact int
}

// ID returns the unique ID of the unit in the system
func (u Unit) ID() string {
	return u.String()
}

// String returns the unit its name and its factor between parentheses
func (u Unit) String() string {
	return u.Name + "(" + strconv.Itoa(u.Fact) + ")"
}

// NewUnit creates a new unit from a string, see 'String()'
func NewUnit(s string) Unit {
	s = strings.TrimSuffix(s, ")")
	ss := strings.Split(s, "(")
	u := Unit{"unknown", 0}
	if len(ss) > 0 {
		u.Name = ss[0]
	}
	if len(ss) > 1 {
		u.Fact, _ = strconv.Atoi(ss[1])
	}
	return u
}

// Quant represents a quantity with unit and value (int)
type Quant struct {
	Unit
	Val int
}

// String returns the value and unit of the quantity
func (q Quant) String() string {
	return strconv.Itoa(q.Val) + " " + q.Unit.String()
}

// Total returns the unit factor by the value of the quantity
func (q Quant) Total() int {
	return q.Fact * q.Val
}

// Empty creates a quantity based on q
func (q Quant) Empty() Quant {
	return Quant{q.Unit, 0}
}

// NewAdd creates a quantity based on q
// and setsg its value to the sum of q value and add
func (q Quant) NewAdd(add int) Quant {
	return Quant{q.Unit, q.Val + add}
}

// NewSub creates a quantity based on q
// and sets its value to the subtraction of sub to q value
func (q Quant) NewSub(sub int) Quant {
	return Quant{q.Unit, q.Val - sub}
}

// NewSet creates a quantity based on q
// and sets its value to set
func (q Quant) NewSet(set int) Quant {
	return Quant{q.Unit, set}
}

// Add creates a quantity based on q1
// and sets its value to the sum of q1 value and q2 value
func Add(q1, q2 Quant) Quant {
	return Quant{q1.Unit, q1.Val + q2.Val}
}

// Sub creates a quantity based on q2
// and sets its value to the subtraction of q2 value to q1 value
func Sub(q1, q2 Quant) Quant {
	return Quant{q1.Unit, q1.Val - q2.Val}
}

// MapToSlice returns a slice quantities from a map of  quantities
func MapToSlice(m map[string]Quant) []Quant {
	var out []Quant
	for _, v := range m {
		out = append(out, v)
	}
	return out
}

// SliceTotal returns the sum of all quantities total
func SliceTotal(qs []Quant) int {
	var total int
	for _, v := range qs {
		total += v.Total()
	}
	return total
}

// TrimSliceOnTotal trims quantities until the sum of their total reach lim
func TrimSliceOnTotal(qs []Quant, lim int) []Quant {
	var out []Quant
	for _, v := range qs {
		out = append(out, v)
		if SliceTotal(out) >= lim {
			return out
		}
	}
	return out
}

// CopyMap copies a map of quantities
func CopyMap(m map[string]Quant) map[string]Quant {
	n := make(map[string]Quant, 0)
	for k, v := range m {
		n[k] = v
	}
	return n
}

// ByFactDesc sorts quantities by descending unit factor
type ByFactDesc []Quant

func (f ByFactDesc) Len() int           { return len(f) }
func (f ByFactDesc) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }
func (f ByFactDesc) Less(i, j int) bool { return f[i].Fact > f[j].Fact }

// ByFactAsc sorts quantities by ascending unit factor
type ByFactAsc []Quant

func (f ByFactAsc) Len() int           { return len(f) }
func (f ByFactAsc) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }
func (f ByFactAsc) Less(i, j int) bool { return f[i].Fact < f[j].Fact }
