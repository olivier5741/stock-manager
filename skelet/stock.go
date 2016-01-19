package skelet

import (
	. "github.com/olivier5741/stock-manager/item"
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
	Items
	Date string
}

// code duplication !! Initialization is different when type is composed
func (i InCmd) Id() string {
	return i.StockName
}

func (i OutCmd) Id() string {
	return i.StockName
}

func (i InventoryCmd) Id() string {
	return i.StockName
}

type InSubmitted struct {
	StockEvent
	In, Stock Items
}

type OutSubmitted struct {
	StockEvent
	Out, Stock Items
}

type InventorySubmitted struct {
	StockEvent
	Inventory, Stock Items
}

type StockEvent struct {
	StockName string
	Date      time.Time
}
