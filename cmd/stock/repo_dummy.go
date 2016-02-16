package stock

import (
	"github.com/olivier5741/stock-manager/skelet"
	"github.com/olivier5741/stock-manager/stock"
)

type DummyStockRepo struct {
	skelet.DumEvtRepo
}

func MakeDummyStockRepo() *DummyStockRepo {
	d := DummyStockRepo{*skelet.MakeDumEvtRepo(stock.MakeStock, stock.FromActions)}
	return &d
}
