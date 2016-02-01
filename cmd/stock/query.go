package stock

import (
	"github.com/olivier5741/stock-manager/item"
	sk "github.com/olivier5741/stock-manager/skelet"
	stockM "github.com/olivier5741/stock-manager/stock/main"
)

func (endPt EndPt) ProdValEvolution(id string) (data map[string]item.Items) {
	// TODO : should be generated when even arrives
	acts := endPt.Db.GetAllEvents(id)
	data = make(map[string]item.Items, len(acts))
	state := make(item.Items,0)

	i := 0
loop:
	for _, a := range acts {
		switch act := a.(sk.Event).Act.(type) {
		case stockM.In:
			state.Add(act.Items)
		case stockM.Out:
			state.Sub(act.Items)
		case stockM.Inventory:
			state = act.Items.Copy()
		default:
			continue loop
		}
		data[a.(sk.Event).Date] = state.Copy()
		i++
	}

	return

}
