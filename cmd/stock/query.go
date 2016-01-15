package stock

import (
	. "github.com/olivier5741/stock-manager/stock"
	"sort"
)

func (endPt EndPt) StocksQuery() (s []*Stock, err error) {
	r, err := endPt.Db.GetAll()
	if err != nil {
		return nil, err
	}

	for _, i := range r {
		s = append(s, i.(*Stock))
	}
	return
}

func MakeProdInStockTable() *ProdInStockTable {
	return &ProdInStockTable{Table: make(map[string]ProdInStockLine)}
}

type ProdInStockLine struct {
	Prod string
	Vals []string
}

type ProdInStockTable struct {
	Stocks []string
	Table  map[string]ProdInStockLine
}

// I could put some computing inside methods ...
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

func (p ProdInStockTable) ToProductStringLines() (lines [][]string) {

	lines = [][]string{append([]string{"product"}, p.Stocks...)}

	keys := sort.StringSlice{}
	for key, _ := range p.Table {
		keys = append(keys, key)
	}

	// Le sort du string dit que les majuscules sont toujours plus grande que les minuscules
	keys.Sort()

	for _, key := range keys {
		item := p.Table[key]
		line := []string{item.Prod}
		line = append(line, item.Vals...)
		lines = append(lines, line)
	}

	return
}
