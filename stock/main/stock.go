package stock

import (
	"github.com/olivier5741/stock-manager/item/items"
	. "github.com/olivier5741/stock-manager/skelet"
	. "github.com/olivier5741/stock-manager/stock/skelet"
)

type Stock struct {
	Name string
	items.Items
}

func MakeStock(name string) Ider {
	return &Stock{name, items.Items{}}
}

func FromActions(acts []interface{}, id string) Ider {
	stock := MakeStock(id).(*Stock)
	for _, act := range acts {
		switch act := act.(type) {
		case In:
			stock.Items = items.Add(stock.Items, act.Items)
		case Out:
			stock.Items = items.Sub(stock.Items, act.Items)
		case Inventory:
			stock.Items = act.Items.Copy()
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
	items.Items
}

func (s Stock) ID() string {
	return s.Name
}

func (s *Stock) SubmitIn(i InCmd) (e In, err error) {
	s.Items = items.Add(s.Items, i.Items)
	e = In{i.Items}
	return
}

func (s *Stock) SubmitOut(o OutCmd) (e Out, err error) {
	s.Items = items.Sub(s.Items, o.Items)
	e = Out{o.Items}
	return
}

func (s *Stock) SubmitInventory(i InventoryCmd) (e Inventory, err error) {
	s.Items = i.Items
	e = Inventory{i.Items}
	return
}

func (s *Stock) RenameStock(r RenameCmd) (e Rename, err error) {
	s.Name = r.Name
	e = Rename{r.Name}
	return
}
