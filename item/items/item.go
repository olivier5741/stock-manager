// Package items provide addition, subtraction as well as
// copy, missing, empty actions to a list of item
package items

import (
	"github.com/olivier5741/stock-manager/item/amount"
	"strconv"
)

// T a list of items
type T map[string]Item

// Item with related product and the value
type Item struct {
	Prod   Prod
	Amount amount.A
}

// String print the product and value of the item
func (it Item) String() string {
	return it.Prod.String() + ": " + it.Amount.String()
}

// Prod the (item) product
type Prod string

// String print the product
func (p Prod) String() string {
	return string(p)
}

// Add creates a list by adding value of each matching item
// and then apprend the non-matching one
func Add(ins, adds T) T {
	its := ins.Copy()
	for key, add := range adds {
		if it, ok := its[key]; ok {
			add.Amount = amount.Add(it.Amount, add.Amount)
		}
		its[key] = add
	}
	return its
}

// Sub creates a list by subtracting value from right item
// to matching left item, append the non-matching left items
// and append the negative of the right non-matching items
func Sub(ins, subs T) T {
	its := ins.Copy()
	for key, sub := range subs {
		if it, ok := its[key]; ok {
			sub.Amount = amount.Sub(it.Amount, sub.Amount)
		} else {
			sub.Amount = sub.Amount.Neg()
		}
		its[key] = sub
	}
	return its
}

// Missing creates a list by comparing what is missing in
// the right list : non-existant items from the left list and
// the difference between matching items
// if bigger than 0 or if cannot be compiled as an integer
func Missing(its, exps T) T {
	out := map[string]Item{}
	for key, exp := range exps {
		if it, ok := its[key]; ok {
			if diff, no, intDiff := amount.Diff(exp.Amount, it.Amount); no && intDiff > 0 {
				out[key] = Item{it.Prod, diff}
			}
		} else {
			out[key] = Item{exp.Prod, exp.Amount.Copy()}
		}
	}
	return out
}

// Empty creates a list by setting the value of each item to empty
func (its T) Empty() T {
	out := map[string]Item{}
	for key, it := range its {
		out[key] = Item{it.Prod, it.Amount.Empty()}
	}
	return out
}

// Copy creates a list by copying each item
func (its T) Copy() T {
	out := make(T, len(its))
	for key, orig := range its {
		out[key] = orig
	}
	return out
}

// should call number of units on val ...
func (its T) MaxUnit() int {
	var max int
	for _, it := range its {
		// TODO method for it.Val.Quants
		if len(it.Amount.Quants()) > max {
			max = len(it.Amount.Quants())
		}
	}
	return max
}

func (its T) StringSlice() [][]string {
	var out [][]string
	max := its.MaxUnit()
	for _, it := range its {
		out = append(out, it.StringSlice(max))
	}
	return out
}

func FromSlice(items []Item) T {
	out := T{}
	for _, item := range items {
		out[string(item.Prod)] = item
	}
	return out
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

func ItemsMapToStringMapTable(itsmap map[string]T) map[string]map[string]string {
	out := make(map[string]map[string]string, 0)
	for date, its := range itsmap {
		newRow := make(map[string]string)
		for prod, it := range its {
			newRow[prod] = strconv.Itoa(it.Amount.TotalWith())
		}
		out[date] = newRow
	}
	return out
}
