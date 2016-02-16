package stock

import (
	"github.com/olivier5741/stock-manager/item/items"
	"github.com/olivier5741/stock-manager/skelet"
)

type Stock struct {
	Name string
	items.I
}

func MakeStock(name string) skelet.Ider {
	return &Stock{name, items.I{}}
}

func (s Stock) ID() string {
	return s.Name
}

func FromActions(acts []interface{}, id string) skelet.Ider {
	s := MakeStock(id).(*Stock)
	for _, act := range acts {
		switch act := act.(type) {
		case In:
			s.I = items.Add(s.I, act.I)
		case Out:
			s.I = items.Sub(s.I, act.I)
		case Inventory:
			s.I = act.I.Copy()
		case Rename:
			s.Name = act.Name
		}
	}
	return s
}

func (s *Stock) SubmitIn(i InCmd) (e In, err error) {
	s.I = items.Add(s.I, i.I)
	e = In{i.I}
	return
}

func (s *Stock) SubmitOut(o OutCmd) (e Out, err error) {
	s.I = items.Sub(s.I, o.I)
	e = Out{o.I}
	return
}

func (s *Stock) SubmitInventory(i InventoryCmd) (e Inventory, err error) {
	s.I = i.I
	e = Inventory{i.I}
	return
}

func (s *Stock) RenameStock(r RenameCmd) (e Rename, err error) {
	s.Name = r.Name
	e = Rename{r.Name}
	return
}


type InCmd ItemsCmd
type OutCmd ItemsCmd
type InventoryCmd ItemsCmd
type RenameCmd struct {
	Name string
}

type ItemsCmd struct {
	StockName string
	items.I
	Date string
}

// code duplication !! Initialization is different when type is composed
func (i InCmd) ID() string {
	return i.StockName
}

func (i OutCmd) ID() string {
	return i.StockName
}

func (i InventoryCmd) ID() string {
	return i.StockName
}

type Rename struct {
	Name string
}

type In ItemsAction
type Out ItemsAction
type Inventory ItemsAction

type ItemsAction struct {
	items.I
}