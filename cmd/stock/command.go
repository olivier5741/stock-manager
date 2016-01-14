package stock

import (
	. "github.com/olivier5741/stock-manager/item"
	. "github.com/olivier5741/stock-manager/skelet"
	. "github.com/olivier5741/stock-manager/stock"
	"time"
)

func (endPt *EndPt) HandleIn(agg interface{}, cmd interface{}) (event interface{}, extEvent interface{}, err error) {
	stock := agg.(*Stock)
	cmdIn := cmd.(InCmd)

	in, err := stock.SubmitIn(cmdIn)
	if err != nil {
		return nil, nil, err
	}

	return in, InSubmitted{
		StockEvent{in.StockName, time.Now()},
		in.Items, stock.Items.Copy()}, nil
}

func (endPt *EndPt) HandleOut(agg interface{}, cmd interface{}) (event interface{}, extEvent interface{}, err error) {
	stock := agg.(*Stock)
	cmdOut := cmd.(OutCmd)

	out, err := stock.SubmitOut(cmdOut)
	if err != nil {
		return nil, nil, err
	}

	return out, OutSubmitted{
		StockEvent{out.StockName, time.Now()},
		out.Items, stock.Items.Copy()}, nil
}

func (endPt *EndPt) HandleInventory(agg interface{}, cmd interface{}) (event interface{}, extEvent interface{}, err error) {
	stock := agg.(*Stock)
	cmdInv := cmd.(InventoryCmd)

	inv, err := stock.SubmitInventory(cmdInv)
	if err != nil {
		return nil, nil, err
	}

	return inv, InventorySubmitted{
		StockEvent{inv.StockName, time.Now()},
		inv.Items, stock.Items.Copy()}, nil
}
