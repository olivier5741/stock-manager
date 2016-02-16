package stock

import (
	"github.com/olivier5741/stock-manager/skelet"
	"github.com/olivier5741/stock-manager/stock"
	"github.com/olivier5741/stock-manager/stock/stockevent"
	"time"
)

func (endPt *EndPt) HandleIn(agg interface{}, cmd interface{}) (event skelet.Event, extEvent interface{}, err error) {
	s := agg.(*stock.Stock)
	cmdIn := cmd.(stock.InCmd)

	in, err := s.SubmitIn(cmdIn)
	if err != nil {
		return skelet.Event{}, nil, err
	}

	return skelet.Event{cmdIn.Date, in}, stockevent.InSubmitted{
		stockevent.StockEvent{s.Name, time.Now()},
		in.I, s.I}, nil
}

func (endPt *EndPt) HandleOut(agg interface{}, cmd interface{}) (event skelet.Event, extEvent interface{}, err error) {
	s := agg.(*stock.Stock)
	cmdOut := cmd.(stock.OutCmd)

	out, err := s.SubmitOut(cmdOut)
	if err != nil {
		return skelet.Event{}, nil, err
	}

	return skelet.Event{cmdOut.Date, out}, stockevent.OutSubmitted{
		stockevent.StockEvent{s.Name, time.Now()},
		s.I, s.I}, nil
}

func (endPt *EndPt) HandleInventory(agg interface{}, cmd interface{}) (event skelet.Event, extEvent interface{}, err error) {
	s := agg.(*stock.Stock)
	cmdInv := cmd.(stock.InventoryCmd)

	inv, err := s.SubmitInventory(cmdInv)
	if err != nil {
		return skelet.Event{}, nil, err
	}

	return skelet.Event{cmdInv.Date, inv}, stockevent.InventorySubmitted{
		stockevent.StockEvent{s.Name, time.Now()},
		inv.I, s.I}, nil
}
