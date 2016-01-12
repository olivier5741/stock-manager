package order

import (
	. "github.com/olivier5741/stock-manager/item"
)

type Order struct {
	Items Items
}

type StockState struct {
	Stocks map[string]Items
}

func (s StockState) WhatIsMissing(min Items) Items {
	total := make(Items)
	for _, stock := range s.Stocks {
		total.Add(stock)
	}
	return total.Missing(min)
}

func (s StockState) StockUpdate(id string, its Items) {
	s.Stocks[id] = its
}
