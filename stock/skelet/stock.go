package stock

import (
	. "github.com/olivier5741/stock-manager/item/items"
	"time"
)

type InCmd ItemsCmd
type OutCmd ItemsCmd
type InventoryCmd ItemsCmd
type RenameCmd struct {
	Name string
}

type ItemsCmd struct {
	StockName string
	T
	Date string
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

type InSubmitted struct {
	StockEvent
	In, Stock T
}

type OutSubmitted struct {
	StockEvent
	Out, Stock T
}

type InventorySubmitted struct {
	StockEvent
	Inventory, Stock T
}

type StockEvent struct {
	StockName string
	Date      time.Time
}
