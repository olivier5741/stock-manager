package stock

import (
	"github.com/olivier5741/stock-manager/item/items"
	"github.com/olivier5741/stock-manager/skelet"
	"github.com/olivier5741/stock-manager/stock"
)

func (endPt EndPt) ProdValEvolution(id string) (data map[string]items.I) {
	// TODO : should be generated when even arrives
	acts := endPt.Db.GetAllEvents(id)
	data = make(map[string]items.I, len(acts))
	state := make(items.I, 0)

	i := 0
loop:
	for _, a := range acts {
		switch act := a.(skelet.Event).Act.(type) {
		case stock.In:
			state = items.Add(state, act.I)
		case stock.Out:
			state = items.Sub(state, act.I)
		case stock.Inventory:
			state = act.I.Copy()
		default:
			continue loop
		}
		data[a.(skelet.Event).Date] = state.Copy()
		i++
	}

	return

}
