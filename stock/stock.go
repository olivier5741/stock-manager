package stock

import (
	log "github.com/Sirupsen/logrus"
	. "github.com/olivier5741/stock-manager/item"
	. "github.com/olivier5741/stock-manager/skelet"
)

type Stock struct {
	Name string
	Items
}

func MakeStock(name string) Ider {
	return &Stock{name, Items{}}
}

func FromActions(acts []interface{}, id string) Ider {
	stock := MakeStock(id).(*Stock)
	log.Debug("EVENTS TO REPLAY")
	log.Debug(acts)
	for _, act := range acts {
		switch act := act.(type) {
		case In:
			stock.Items.Add(act.Items)
		case Out:
			stock.Items.Sub(act.Items)
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
	Items
}

func (s Stock) Id() string {
	return s.Name
}

func (s *Stock) SubmitIn(i InCmd) (e In, err error) {
	s.Items.Add(i.Items)
	e = In{i.Items}
	return
}

func (s *Stock) SubmitOut(o OutCmd) (e Out, err error) {
	log.Debug("ITEMS")
	log.Debug(s.Items)
	log.Debug(o.Items)
	s.Items.Sub(o.Items)
	log.Debug("ITEMS after sub")
	log.Debug(s.Items)
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
