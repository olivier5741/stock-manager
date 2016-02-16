package item

import (
	"github.com/olivier5741/stock-manager/item/amount"
	"strconv"
)

// Item with related product and the value
type Item struct {
	Prod   Prod
	Amount amount.A
}

// String print the product and value of the item
func (it Item) String() string {
	return it.Prod.String() + ": " + it.Amount.String()
}

// remove limit to somewhere else ...
func (it Item) StringSlice(unitNb int) []string {
	s := make([]string, unitNb*2+1)
	s[0] = it.Prod.String()
	count := 1
	for _, u := range it.Amount.ValsWithByFactDesc() {
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