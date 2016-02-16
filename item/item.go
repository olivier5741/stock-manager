package item

import (
	"strconv"
	"github.com/olivier5741/stock-manager/item/amount"
)

// I Item  with related product and the value
type I struct {
	Prod   Prod
	Amount amount.A
}

// String print the product and value of the item
func (it I) String() string {
	return it.Prod.String() + ": " + it.Amount.String()
}

// StringSlice formats the item as a csv row (string array).
// Product; first value; first unit, second value; second unit; ...
func (it I) StringSlice(unitNb int) []string {
	s := make([]string, unitNb*2+1)
	s[0] = it.Prod.String()
	count := 1
	for _, u := range it.Amount.QuantsWithByFactDesc() {
		if count == unitNb*2+1 {
			break
		}
		s[count] = strconv.Itoa(u.Val)
		s[count+1] = u.Unit.String()
		count += 2
	}
	return s
}

// Prod the (item) product
type Prod string

// String print the product
func (p Prod) String() string {
	return string(p)
}
