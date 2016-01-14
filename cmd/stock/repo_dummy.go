package stock

import (
	. "github.com/olivier5741/stock-manager/skelet"
	. "github.com/olivier5741/stock-manager/stock"
	"log"
)

type DummyStockRepository struct {
	db map[string][]interface{}
}

func MakeDummyStockRepository() *DummyStockRepository {
	db := map[string][]interface{}{}
	return &DummyStockRepository{db}
}

func (r DummyStockRepository) GetAll() (aggs []interface{}) {

	log.Println(r)
	for key, _ := range r.db {
		out, _ := r.Get(key)
		aggs = append(aggs, out)
	}
	return
}

func (r DummyStockRepository) Get(id string) (s Ider, err error) {
	events, ok := r.db[id]
	if !ok {
		return MakeStock(id), nil
	}
	return FromActions(events), nil
}

func (r *DummyStockRepository) Put(id string, event interface{}) (err error) {
	r.db[id] = append(r.db[id], event)
	//log.Println("Database State")
	//log.Println(r.db)
	return
}
