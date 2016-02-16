package quant

import (
	"strconv"
)

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
