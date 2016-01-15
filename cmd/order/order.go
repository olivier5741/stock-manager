package order

import (
	. "github.com/olivier5741/stock-manager/order"
	. "github.com/olivier5741/stock-manager/skelet"
)

// I deleted Unique persister interface...
type EndPt struct {
	db UniquePersister
}

type StockStateDatabase StockState

func (s StockStateDatabase) Get() interface{} {
	return StockState(s)
}

func (s StockStateDatabase) Put(i interface{}) {
	s = StockStateDatabase(i.(StockState))
}

func (endPt EndPt) HandleInSubmitted(i InSubmitted) {
	ss := endPt.db.Get().(StockState)
	ss.StockUpdate(i.StockName, i.Stock.Copy())
	endPt.db.Put(ss)
}
