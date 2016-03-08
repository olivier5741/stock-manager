package stock

import (
	"github.com/olivier5741/stock-manager/item/items"
	"github.com/olivier5741/stock-manager/skelet"
)

type Stock struct {
	Name string
	items.Items
	Min items.Items
}

func ProdsUpdateFromStringTable(t [][]string) (mins items.Items, units items.Items) {
	var minsTable, unitsTable [][]string

	for _,row := range t {
		minsTable = append(minsTable,[]string{row[0],row[1],row[2]})
		units := []string{row[0]}
		for _,u := range row[3:]{
			units = append(units,"0",u)
		}
		unitsTable = append(unitsTable,units)
	}

	return items.FromStringTable(minsTable),
		items.FromStringTable(unitsTable)
	
}

func MakeStock(name string) skelet.Ider {
	return &Stock{name, items.Items{}, items.Items{}}
}

func (s Stock) ID() string {
	return s.Name
}

func FromActions(acts []interface{}, id string) skelet.Ider {
	s := MakeStock(id).(*Stock)
	for _, act := range acts {
		switch act := act.(type) {
		case In:
			s.Items = items.Add(s.Items, act.Items)
		case Out:
			s.Items = items.Sub(s.Items, act.Items)
		case Inventory:
			s.Items = act.Items.Copy()
		case ProdsUpdate:
			s.Min = act.Mins.Copy()
			s.Items = items.Add(s.Items,act.Units)
		case Rename:
			s.Name = act.Name
		}
	}
	return s
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

func (s *Stock) UpdateProds(p ProdsUpdateCmd) (e ProdsUpdate, err error) {
	s.Min = p.Mins
	s.Items = items.Add(s.Items,p.Units)
	e = ProdsUpdate{p.Mins,p.Units}
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

type ItemsCmd struct {
	StockName string
	items.Items
	Date string
}

type ProdsUpdateCmd struct {
	StockName string
	Mins items.Items
	Units items.Items
	Date string
}

type RenameCmd struct {
	Name string
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

func (i ProdsUpdateCmd) ID() string {
	return i.StockName
}

type Rename struct {
	Name string
}

type In ItemsAct
type Out ItemsAct
type Inventory ItemsAct
type ProdsUpdate struct {
	Mins items.Items
	Units items.Items // should be of another type
}

type ItemsAct struct {
	items.Items
}
