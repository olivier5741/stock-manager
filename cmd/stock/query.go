package stock

import (
	. "github.com/olivier5741/stock-manager/item"
	. "github.com/olivier5741/stock-manager/skelet"
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

func (endPt EndPt) ProdValEvolution(id string) (data map[string]Items) {
	acts := endPt.Db.GetAllEvents(id)
	data = make(map[string]Items, len(acts))
	state := make(Items, 0)

	i := 0
loop:
	for _, a := range acts {
		switch act := a.(Event).Act.(type) {
		case In:
			state.Add(act.Items)
		case Out:
			state.Sub(act.Items)
		case Inventory:
			state = act.Items.Copy()
		default:
			continue loop
		}
		data[a.(Event).Date] = state.Copy()
		i++
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

func ToProductStringLines(items Items) (lines [][]string) {

	lines = [][]string{[]string{"Prod", "Val"}}

	keys := sort.StringSlice{}
	for key, _ := range items {
		keys = append(keys, key)
	}

	// Le sort du string dit que les majuscules sont toujours plus grande que les minuscules
	keys.Sort()

	for _, key := range keys {
		item := items[key]
		line := []string{item.Prod.String(), item.Val.String()}
		lines = append(lines, line)
	}

	return
}
