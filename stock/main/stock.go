package stock

import (
	"github.com/olivier5741/stock-manager/item/items"
	. "github.com/olivier5741/stock-manager/skelet"
	. "github.com/olivier5741/stock-manager/stock/skelet"
)

type Stock struct {
	Name string
	items.T
}

func MakeStock(name string) Ider {
	return &Stock{name, items.T{}}
}

func FromActions(acts []interface{}, id string) Ider {
	stock := MakeStock(id).(*Stock)
	for _, act := range acts {
		switch act := act.(type) {
		case In:
			stock.T = items.Add(stock.T, act.T)
		case Out:
			stock.T = items.Sub(stock.T, act.T)
		case Inventory:
			stock.T = act.T.Copy()
		case Rename:
			stock.Name = act.Name
		}
	}
	return stock
}

type Rename struct {
	Name string
}

type In ItemsAction
type Out ItemsAction
type Inventory ItemsAction

type ItemsAction struct {
	items.T
}

func (s Stock) ID() string {
	return s.Name
}

func (s *Stock) SubmitIn(i InCmd) (e In, err error) {
	s.T = items.Add(s.T, i.T)
	e = In{i.T}
	return
}

func (s *Stock) SubmitOut(o OutCmd) (e Out, err error) {
	s.T = items.Sub(s.T, o.T)
	e = Out{o.T}
	return
}

func (s *Stock) SubmitInventory(i InventoryCmd) (e Inventory, err error) {
	s.T = i.T
	e = Inventory{i.T}
	return
}

func (s *Stock) RenameStock(r RenameCmd) (e Rename, err error) {
	s.Name = r.Name
	e = Rename{r.Name}
	return
}
