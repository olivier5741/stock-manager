// Package items provide addition, subtraction as well as
// copy, missing, empty actions to a list of item
package items

import (
	"github.com/olivier5741/stock-manager/item/val"
	"strconv"
)

// Items a list of items
type Items map[string]Item

// Item with related product and the value
type Item struct {
	Prod Prod
	Val  val.Val
}

// String print the product and value of the item
func (i Item) String() string {
	return i.Prod.String() + ": " + i.Val.String()
}

// Prod the (item) product
type Prod string

// String print the product
func (p Prod) String() string {
	return string(p)
}

// Add creates a list by adding value of each matching item
// and then apprend the non-matching one
func Add(ins, adds Items) Items {
	its := ins.Copy()
	for key, add := range adds {
		if it, ok := its[key]; ok {
			add.Val = val.Add(it.Val, add.Val)
		}
		its[key] = add
	}
	return its
}

// Add creates a list by subtracting value from right item
// to matching left item, append the non-matching left items
// and append the negative of the right non-matching items
func Sub(ins, subs Items) Items {
	its := ins.Copy()
	for key, sub := range subs {
		if it, ok := its[key]; ok {
			sub.Val = val.Sub(it.Val, sub.Val)
		} else {
			sub.Val = sub.Val.Neg()
		}
		its[key] = sub
	}
	return its
}

// Missing creates a list by comparing what is missing in
// the right list : non-existant items from the left list and
// the difference between matching items
// if bigger than 0 or if cannot be compiled as an integer
func Missing(its, exps Items) (out Items) {
	out = map[string]Item{}
	for key, exp := range exps {
		if it, ok := its[key]; ok {
			if diff, no, intDiff := val.Diff(exp.Val, it.Val); no && intDiff > 0 {
				out[key] = Item{it.Prod, diff}
			}
		} else {
			out[key] = Item{exp.Prod, exp.Val.Copy()}
		}
	}
	return
}

// Empty creates a list by setting the value of each item to empty
func (its Items) Empty() (out Items) {
	out = map[string]Item{}
	for key, it := range its {
		out[key] = Item{it.Prod, it.Val.Empty()}
	}
	return
}

// Copy creates a list by copying each item
func (its Items) Copy() (out Items) {
	out = make(Items, len(its))
	for key, orig := range its {
		out[key] = orig
	}
	return
}

// should call number of units on val ...
func (its Items) MaxUnit() (max int) {
	for _, it := range its {
		if len(it.Val.Vals) > max {
			max = len(it.Val.Vals)
		}
	}
	return
}

func (its Items) StringSlice() [][]string {
	var out [][]string
	max := its.MaxUnit()
	for _, it := range its {
		out = append(out, it.StringSlice(max))
	}
	return out
}

func FromSlice(items []Item) (out Items) {
	out = Items{}
	for _, item := range items {
		out[string(item.Prod)] = item
	}
	return
}

// remove limit to somewhere else ...
func (v Item) StringSlice(unitNb int) []string {
	s := make([]string, unitNb*2+1)
	s[0] = v.Prod.String()
	count := 1
	for _, u := range v.Val.ValsWithByFactorDesc() {
		if count == unitNb*2+1 {
			break
		}
		s[count] = strconv.Itoa(u.Val)
		s[count+1] = u.Unit.String()
		count += 2
	}
	return s
}

func MapItemsMap(its map[string]Items) map[string]map[string]string {
	out := make(map[string]map[string]string, 0)
	for date, it := range its {
		newRow := make(map[string]string)
		for prod, val := range it {
			newRow[prod] = strconv.Itoa(val.Val.TotalWith())
		}
		out[date] = newRow
	}
	return out
}
