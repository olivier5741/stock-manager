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

// QFloat represents a quantity as Q but with a float value
type QFloat struct {
	Unit
	Val float64
}

// String returns the value of qf (2 decimals) and its unit
func (qf QFloat) String() string {
	return strconv.FormatFloat(qf.Val, 'f', 2, 64) +
		" " + qf.Unit.String()
}
