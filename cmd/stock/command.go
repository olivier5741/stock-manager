package stock

import (
	. "github.com/olivier5741/stock-manager/item"
	. "github.com/olivier5741/stock-manager/skelet"
	. "github.com/olivier5741/stock-manager/stock"
	"time"
)

func (endPt *EndPt) HandleIn(agg interface{}, cmd interface{}) (event Event, extEvent interface{}, err error) {
	stock := agg.(*Stock)
	cmdIn := cmd.(InCmd)

	in, err := stock.SubmitIn(cmdIn)
	if err != nil {
		return Event{}, nil, err
	}

	return Event{cmdIn.Date, in}, InSubmitted{
		StockEvent{stock.Name, time.Now()},
		in.Items, stock.Items}, nil
}

func (endPt *EndPt) HandleOut(agg interface{}, cmd interface{}) (event Event, extEvent interface{}, err error) {
	stock := agg.(*Stock)
	cmdOut := cmd.(OutCmd)

	out, err := stock.SubmitOut(cmdOut)
	if err != nil {
		return Event{}, nil, err
	}

	return Event{cmdOut.Date, out}, OutSubmitted{
		StockEvent{stock.Name, time.Now()},
		out.Items, stock.Items}, nil
}

func (endPt *EndPt) HandleInventory(agg interface{}, cmd interface{}) (event Event, extEvent interface{}, err error) {
	stock := agg.(*Stock)
	cmdInv := cmd.(InventoryCmd)

	inv, err := stock.SubmitInventory(cmdInv)
	if err != nil {
		return Event{}, nil, err
	}

	return Event{cmdInv.Date, inv}, InventorySubmitted{
		StockEvent{stock.Name, time.Now()},
		inv.Items, stock.Items}, nil
}
