package stockevent

import (
	"github.com/olivier5741/stock-manager/item/items"
	"time"
)

type InSubmitted struct {
	StockEvent
	In, Stock items.I
}

type OutSubmitted struct {
	StockEvent
	Out, Stock items.I
}

type InventorySubmitted struct {
	StockEvent
	Inventory, Stock items.I
}

type StockEvent struct {
	StockName string
	Date      time.Time
}
