package stock

import (
	"github.com/olivier5741/stock-manager/skelet"
	"github.com/olivier5741/stock-manager/stock/main"
)

type DummyStockRepository struct {
	skelet.DumEvtRepo
}

func MakeDummyStockRepository() *DummyStockRepository {
	d := DummyStockRepository{*skelet.MakeDumEvtRepo(stock.MakeStock, stock.FromActions)}
	return &d
}
