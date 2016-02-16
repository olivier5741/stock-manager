package stockevent

import (
	"github.com/olivier5741/stock-manager/item/items"
	"time"
)

type InSubmitted struct {
	StockEvt
	In, Stock items.Items
}

type OutSubmitted struct {
	StockEvt
	Out, Stock items.Items
}

type InventorySubmitted struct {
	StockEvt
	Inventory, Stock items.Items
}

type StockEvt struct {
	StockName string
	Date      time.Time
}
