package unitval

import (
	"strconv"
	"strings"
)

type Unit struct {
	Name string
	Fact int
}

func (u Unit) ID() string {
	return u.String()
}

func (u Unit) String() string {
	return u.Name + "(" + strconv.Itoa(u.Fact) + ")"
}

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

type TFloat struct {
	Unit
	Val float64
}

func (u TFloat) String() string {
	return strconv.FormatFloat(u.Val, 'f', 2, 64) + " " + u.Unit.String()
}
