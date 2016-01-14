package stock

import (
	. "github.com/olivier5741/stock-manager/item"
	. "github.com/olivier5741/stock-manager/skelet"
)

type Stock struct {
	Name string
	Items
}

func MakeStock(name string) *Stock {
	return &Stock{name, Items{}}
}

type Rename struct {
	Name string
}

type In ItemsAction
type Out ItemsAction
type Inventory ItemsAction

type ItemsAction struct {
	StockName string
	Items
}

func (s Stock) Id() string {
	return s.Name
}

func FromActions(acts []interface{}) (stock *Stock) {
	stock = &Stock{Items: Items{}}
	for _, act := range acts {
		switch act := act.(type) {
		case In:
			stock.Add(act.Items)
			stock.Name = act.StockName
		case Out:
			stock.Sub(act.Items)
			stock.Name = act.StockName
		case Inventory:
			stock.Items = act.Items.Copy()
			stock.Name = act.StockName
		case Rename:
			stock.Name = act.Name
		}
	}
	return
}

func (s *Stock) SubmitIn(i InCmd) (e In, err error) {
	//log.Println("Stock before in")
	//log.Println(s)
	s.Items.Add(i.Items)
	//log.Println("Stock after in")
	//log.Println(s)
	e = In{i.StockName, i.Items.Copy()}
	return
}

func (s *Stock) SubmitOut(o OutCmd) (e Out, err error) {
	//log.Println("Stock before out")
	//log.Println(s)
	s.Items.Sub(o.Items)
	//log.Println("Stock after out")
	//log.Println(s)
	e = Out{o.StockName, o.Items.Copy()}
	return
}

func (s *Stock) SubmitInventory(i InventoryCmd) (e Inventory, err error) {
	//log.Println("Stock before inventory")
	//log.Println(s)
	s.Items = i.Items.Copy()
	//log.Println("Stock after inventory")
	//log.Println(s)
	e = Inventory{i.StockName, i.Items.Copy()}
	return
}

func (s *Stock) RenameStock(r RenameCmd) (e Rename, err error) {
	s.Name = r.Name
	e = Rename{r.Name}
	return
}
