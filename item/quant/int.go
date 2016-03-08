package quant

import (
	"strings"
	"math/big"
	"fmt"
)

// Unit a unit with a name and a factor (must be a factor of smallest
// unit factor in the system)
//
// A system is for instance all the units a particular product can have
// ibuprofen : pill, tab, box, ...

type Unit struct {
	Name string
	Fact *big.Rat
}

// ID returns the unique ID of the unit in the system
func (u Unit) ID() string {
	return u.String()
}

// String returns the unit its name and its factor between parentheses
func (u Unit) String() string {
	return u.Name + "(" + RatToString(u.Fact) + ")"
}

// NewUnit creates a new unit from a string, see 'String()'
func NewUnit(s string) Unit {
	s = strings.TrimSuffix(s, ")")
	ss := strings.Split(s, "(")
	u := Unit{"unknown", &big.Rat{}}
	if len(ss) > 0 {
		u.Name = ss[0]
	}
	if len(ss) > 1 {
		u.Fact = StringToRat(ss[1])
	}
	return u
}

func NewUnits(list []string) []Unit {
	var out []Unit
	for _,it := range list {
		out = append(out, NewUnit(it))
	}
	return out
}

// Quant represents a quantity with unit and value (int)
type Quant struct {
	Unit
	Val *big.Rat
}

func StringToRat(s string) *big.Rat {
	out := new(big.Rat)
	_, _ = fmt.Sscan(s, out)
	return out
}

func RatToString(r *big.Rat) string {
	if r.IsInt() {
		return r.FloatString(0)
	}else {
		return r.FloatString(3)
	}
}

// String returns the value and unit of the quantity
func (q Quant) String() string {
	return RatToString(q.Val) + " " + q.Unit.String()
}

// Total returns the unit factor by the value of the quantity
func (q Quant) Total() *big.Rat {
	return new(big.Rat).Mul(q.Fact,q.Val)
}

// Empty creates a quantity based on q
func (q Quant) Empty() Quant {
	return Quant{q.Unit, &big.Rat{}}
}

// NewAdd creates a quantity based on q
// and setsg its value to the sum of q value and add
func (q Quant) NewAdd(add *big.Rat) Quant {
	return Quant{q.Unit, new(big.Rat).Add(q.Val,add)}
}

// NewSub creates a quantity based on q
// and sets its value to the subtraction of sub to q value
func (q Quant) NewSub(sub *big.Rat) Quant {
	return Quant{q.Unit, new(big.Rat).Sub(q.Val,sub)}
}

// NewSet creates a quantity based on q
// and sets its value to set
func (q Quant) NewSet(set *big.Rat) Quant {
	return Quant{q.Unit, new(big.Rat).Add(&big.Rat{},set)}
}

// Add creates a quantity based on q1
// and sets its value to the sum of q1 value and q2 value
func Add(q1, q2 Quant) Quant {
	return Quant{q1.Unit, new(big.Rat).Add(q1.Val,q2.Val)}
}

// Sub creates a quantity based on q2
// and sets its value to the subtraction of q2 value to q1 value
func Sub(q1, q2 Quant) Quant {
	return Quant{q1.Unit, new(big.Rat).Sub(q1.Val,q2.Val)}
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
func SliceTotal(qs []Quant) *big.Rat {
	total := &big.Rat{}
	for _, v := range qs {
		total.Add(total,v.Total())
	}
	return total
}

// TrimSliceOnTotal trims quantities until the sum of their total reach lim
func TrimSliceOnTotal(qs []Quant, lim *big.Rat) []Quant {
	var out []Quant
	for _, v := range qs {
		out = append(out, v)
		if SliceTotal(out).Cmp(lim) >= 0 {
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
func (f ByFactDesc) Less(i, j int) bool { return f[i].Fact.Cmp(f[j].Fact) > 0 }

// ByFactAsc sorts quantities by ascending unit factor
type ByFactAsc []Quant

func (f ByFactAsc) Len() int           { return len(f) }
func (f ByFactAsc) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }
func (f ByFactAsc) Less(i, j int) bool { return f[i].Fact.Cmp(f[j].Fact) < 0 }
