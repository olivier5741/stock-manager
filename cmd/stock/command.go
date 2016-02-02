package stock

import (
	sk "github.com/olivier5741/stock-manager/skelet"
	stockM "github.com/olivier5741/stock-manager/stock/main"
	stockSk "github.com/olivier5741/stock-manager/stock/skelet"
	"time"
)

func (endPt *EndPt) HandleIn(agg interface{}, cmd interface{}) (event sk.Event, extEvent interface{}, err error) {
	stock := agg.(*stockM.Stock)
	cmdIn := cmd.(stockSk.InCmd)

	in, err := stock.SubmitIn(cmdIn)
	if err != nil {
		return sk.Event{}, nil, err
	}

	return sk.Event{cmdIn.Date, in}, stockSk.InSubmitted{
		stockSk.StockEvent{stock.Name, time.Now()},
		in.T, stock.T}, nil
}

func (endPt *EndPt) HandleOut(agg interface{}, cmd interface{}) (event sk.Event, extEvent interface{}, err error) {
	stock := agg.(*stockM.Stock)
	cmdOut := cmd.(stockSk.OutCmd)

	out, err := stock.SubmitOut(cmdOut)
	if err != nil {
		return sk.Event{}, nil, err
	}

	return sk.Event{cmdOut.Date, out}, stockSk.OutSubmitted{
		stockSk.StockEvent{stock.Name, time.Now()},
		out.T, stock.T}, nil
}

func (endPt *EndPt) HandleInventory(agg interface{}, cmd interface{}) (event sk.Event, extEvent interface{}, err error) {
	stock := agg.(*stockM.Stock)
	cmdInv := cmd.(stockSk.InventoryCmd)

	inv, err := stock.SubmitInventory(cmdInv)
	if err != nil {
		return sk.Event{}, nil, err
	}

	return sk.Event{cmdInv.Date, inv}, stockSk.InventorySubmitted{
		stockSk.StockEvent{stock.Name, time.Now()},
		inv.T, stock.T}, nil
}
