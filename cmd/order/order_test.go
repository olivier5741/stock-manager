package order

import (
	. "github.com/olivier5741/stock-manager/item"
	. "github.com/olivier5741/stock-manager/item_test_shared"
	. "github.com/olivier5741/stock-manager/order"
	. "github.com/olivier5741/stock-manager/skelet"
	"testing"
	"time"
)

func TestHandleInSubmitted(t *testing.T) {
	endPt := EndPt{db: StockStateDatabase(StockState{Stocks: make(map[string]Items)})}

	endPt.HandleInSubmitted(InSubmitted{
		StockEvent: StockEvent{"Bièvre", time.Now()},
		Stock: Items{
			IsoK: Iso3,
			AspK: Asp5,
		}})

	endPt.HandleInSubmitted(InSubmitted{
		StockEvent: StockEvent{"Libin", time.Now()},
		Stock: Items{
			IsoK: Iso3,
			AspK: Asp1,
		}})

	endPt.HandleInSubmitted(InSubmitted{
		StockEvent: StockEvent{"Libramont", time.Now()},
		Stock: Items{
			IsoK: Iso7,
			AspK: Asp6,
		}})

	min := map[string]Item{AspK: Asp20, IsoK: Iso20}
	cmd := endPt.db.Get().(StockState).WhatIsMissing(min)

	exps := map[string]Item{AspK: Asp8, IsoK: Iso7}
	CheckItemsValueAndExistence(t, cmd, exps, "cmd")
}
