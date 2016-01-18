package item

import (
	"strconv"
)

type Items map[string]Item

type Item struct {
	Prod Prod
	Val  Val
}

func (i Item) String() string {
	return i.Prod.String() + ": " + i.Val.String()
}

type Val struct {
	T int
}

func (v Val) String() string {
	return strconv.Itoa(int(v.T))
}

func (v Val) MarshalCSV() (string, error) {
	return string(v.T), nil
}

func (v *Val) UnmarshalCSV(csv string) error {
	if csv == "" {
		v.T = 0
	} else {
		val, err := strconv.ParseInt(csv, 0, 32)
		if err != nil {
			return err
		}
		v.T = int(val)
	}
	return nil
}

type Prod string

func (p Prod) String() string {
	return string(p)
}

func (its Items) Add(adds Items) {
	for key, add := range adds {
		if it, ok := its[key]; ok {
			add.Val.T = it.Val.T + add.Val.T
		}
		its[key] = add
	}
}

func (its Items) Sub(subs Items) {
	for key, sub := range subs {
		if it, ok := its[key]; ok {
			sub.Val.T = it.Val.T - sub.Val.T
		} else {
			sub.Val.T = -sub.Val.T
		}
		its[key] = sub
	}
}

func (its Items) Missing(exps Items) (out Items) {
	out = map[string]Item{}
	for key, exp := range exps {
		if it, ok := its[key]; ok {
			if diff := exp.Val.T - it.Val.T; diff > 0 {
				out[key] = Item{it.Prod, Val{diff}}
			}
		} else {
			out[key] = Item{exp.Prod, Val{exp.Val.T}}
		}
	}
	return
}

// I don't think I'm using this function
func (origs Items) Copy() (out Items) {
	out = make(Items, len(origs))
	for key, orig := range origs {
		out[key] = orig
	}
	return
}

func (its Items) ToStringLines() (out [][]string) {
	out = make([][]string, len(its))
	line := 0
	for _, i := range its {
		out[line] = []string{i.Prod.String(), i.Val.String()}
		line++
	}
	return
}
