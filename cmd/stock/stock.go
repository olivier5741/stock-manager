package stock

import (
	. "github.com/olivier5741/stock-manager/item"
	. "github.com/olivier5741/stock-manager/skelet"
	. "github.com/olivier5741/stock-manager/stock"
	"log"
	"time"
)

type EndPt struct {
	Db EvtSrcPersister
}

func (endPt EndPt) StocksQuery() (s []Stock) {
	for _, i := range endPt.Db.GetAll() {
		s = append(s, i.(Stock))
	}
	return
}

type ProdInStockLine struct {
	Prod string
	Vals []string
}

type ProdInStockTable struct {
	Stocks []string
	Table  map[string]ProdInStockLine
}

func (p *ProdInStockTable) Parse(stocks []Stock) {
	it := 0
	for _, stock := range stocks {
		p.Stocks = append(p.Stocks, stock.Name)
		items := stock.Items.Copy()
		if it != 0 {
			for key, line := range p.Table {
				id := string(line.Prod)
				if item, ok := items[id]; ok {
					line.Vals = append(line.Vals, item.Val.String())
					delete(items, id)
				} else {
					line.Vals = append(line.Vals, "")
				}
				//log.Println(line)
				p.Table[key] = line
			}
		}
		for _, item := range items {
			vals := make([]string, it)
			vals = append(vals, item.Val.String())
			newLine := ProdInStockLine{string(item.Prod), vals}
			p.Table[string(item.Prod)] = newLine
		}
		it++
	}
}

func (endPt EndPt) HandleIn(cmd interface{}) (err error) {
	i := cmd.(InCmd)
	// Transaction 1
	id := i.StockName
	stock := endPt.Db.Get(id).(Stock)

	in, err := stock.SubmitIn(i)
	if err != nil {
		return
	}

	endPt.Db.Put(id, in)

	// Transaction 2
	inSubEvent := InSubmitted{
		StockEvent{i.StockName, time.Now()},
		in.Items, stock.Items.Copy()}

	log.Println("InSubmittedEvent, stock: " + inSubEvent.StockName)
	//log.Println(inSubEvent)

	return
}

func (endPt EndPt) HandleOut(cmd interface{}) (err error) {
	i := cmd.(OutCmd)
	// Transaction 1
	id := i.StockName
	stock := endPt.Db.Get(id).(Stock)

	out, err := stock.SubmitOut(i)
	if err != nil {
		return
	}

	endPt.Db.Put(id, out)

	// Transaction 2
	outSubEvent := OutSubmitted{
		StockEvent{i.StockName, time.Now()},
		out.Items, stock.Items.Copy()}

	log.Println("OutSubmittedEvent, stock: " + outSubEvent.StockName)
	//log.Println(outSubEvent)

	return
}

func (endPt EndPt) HandleInventory(cmd interface{}) (err error) {
	i := cmd.(InventoryCmd)
	// Transaction 1
	id := i.StockName
	stock := endPt.Db.Get(id).(Stock)

	inv, err := stock.SubmitInventory(i)
	if err != nil {
		return
	}

	endPt.Db.Put(id, inv)

	// Transaction 2
	invSubEvent := InventorySubmitted{
		StockEvent{i.StockName, time.Now()},
		inv.Items, endPt.Db.Get(id).(Stock).Items}
	// résoudre ce problème d'assignation !!

	log.Println("InventorySubmittedEvent, stock: " + invSubEvent.StockName)
	//log.Println(invSubEvent)

	return
}

type DummyStockRepository struct {
	db map[string][]interface{}
}

func MakeDummyStockRepository() DummyStockRepository {
	db := map[string][]interface{}{}
	return DummyStockRepository{db}
}

func (r DummyStockRepository) GetAll() (aggs []interface{}) {
	for key, _ := range r.db {
		aggs = append(aggs, r.Get(key))
	}
	return
}

func (r DummyStockRepository) Get(id string) (s interface{}) {
	events, ok := r.db[id]
	if !ok {
		return MakeStock(id)
	}
	return FromActions(events)
}

func (r DummyStockRepository) Put(id string, event interface{}) {
	r.db[id] = append(r.db[id], event)
	//log.Println("Database State")
	//log.Println(r.db)
}
