package stock

import (
	. "github.com/olivier5741/stock-manager/stock"
)

func (endPt EndPt) StocksQuery() (s []*Stock) {
	for _, i := range endPt.Db.GetAll() {
		s = append(s, i.(*Stock))
	}
	return
}

type ProdInStockLine struct {
	Prod string
	Vals []string
}

type ProdInStockTable struct {
	Stocks []string
	Table  map[string]ProdInStockLine
}

func (p *ProdInStockTable) Parse(stocks []*Stock) {
	it := 0
	for _, stock := range stocks {
		p.Stocks = append(p.Stocks, stock.Name)
		items := stock.Items.Copy()
		if it != 0 {
			for key, line := range p.Table {
				id := string(line.Prod)
				if item, ok := items[id]; ok {
					line.Vals = append(line.Vals, item.Val.String())
					delete(items, id)
				} else {
					line.Vals = append(line.Vals, "")
				}
				//log.Println(line)
				p.Table[key] = line
			}
		}
		for _, item := range items {
			vals := make([]string, it)
			vals = append(vals, item.Val.String())
			newLine := ProdInStockLine{string(item.Prod), vals}
			p.Table[string(item.Prod)] = newLine
		}
		it++
	}
}
