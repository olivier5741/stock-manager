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
		stockevent.StockEvt{s.Name, time.Now()},
		in.Items, s.Items}, nil
}

func (endPt *EndPt) HandleOut(agg interface{}, cmd interface{}) (event skelet.Event, extEvent interface{}, err error) {
	s := agg.(*stock.Stock)
	cmdOut := cmd.(stock.OutCmd)

	out, err := s.SubmitOut(cmdOut)
	if err != nil {
		return skelet.Event{}, nil, err
	}

	return skelet.Event{cmdOut.Date, out}, stockevent.OutSubmitted{
		stockevent.StockEvt{s.Name, time.Now()},
		s.Items, s.Items}, nil
}

func (endPt *EndPt) HandleInventory(agg interface{}, cmd interface{}) (event skelet.Event, extEvent interface{}, err error) {
	s := agg.(*stock.Stock)
	cmdInv := cmd.(stock.InventoryCmd)

	inv, err := s.SubmitInventory(cmdInv)
	if err != nil {
		return skelet.Event{}, nil, err
	}

	return skelet.Event{cmdInv.Date, inv}, stockevent.InventorySubmitted{
		stockevent.StockEvt{s.Name, time.Now()},
		inv.Items, s.Items}, nil
}

func (endPt *EndPt) HandleMinimum(agg interface{}, cmd interface{}) (event skelet.Event, extEvent interface{}, err error) {
	s := agg.(*stock.Stock)
	cmdMin := cmd.(stock.MinimumCmd)

	min, err := s.UpdateMinimum(cmdMin)
	if err != nil {
		return skelet.Event{}, nil, err
	}

	return skelet.Event{cmdMin.Date, min}, stockevent.MinimumUpdated{
		stockevent.StockEvt{s.Name, time.Now()},
		min.Items, s.Items}, nil
}
