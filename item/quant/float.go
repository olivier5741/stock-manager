package quant

import (
	"strconv"
)

// QuantFloat represents a quantity as Q but with a float value
type QuantFloat struct {
	Unit
	Val float64
}

// String returns the value of qf (2 decimals) and its unit
func (qf QuantFloat) String() string {
	return strconv.FormatFloat(qf.Val, 'f', 2, 64) +
		" " + qf.Unit.String()
}
