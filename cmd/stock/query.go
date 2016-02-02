package stock

import (
	"github.com/olivier5741/stock-manager/item/items"
	sk "github.com/olivier5741/stock-manager/skelet"
	stockM "github.com/olivier5741/stock-manager/stock/main"
)

func (endPt EndPt) ProdValEvolution(id string) (data map[string]items.T) {
	// TODO : should be generated when even arrives
	acts := endPt.Db.GetAllEvents(id)
	data = make(map[string]items.T, len(acts))
	state := make(items.T, 0)

	i := 0
loop:
	for _, a := range acts {
		switch act := a.(sk.Event).Act.(type) {
		case stockM.In:
			state = items.Add(state, act.T)
		case stockM.Out:
			state = items.Sub(state, act.T)
		case stockM.Inventory:
			state = act.T.Copy()
		default:
			continue loop
		}
		data[a.(sk.Event).Date] = state.Copy()
		i++
	}

	return

}
