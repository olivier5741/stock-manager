package stock

import (
	. "github.com/olivier5741/stock-manager/item"
	. "github.com/olivier5741/stock-manager/skelet"
	. "github.com/olivier5741/stock-manager/stock"
)

func (endPt EndPt) ProdValEvolution(id string) (data map[string]Items) {
	// TODO : should be generated when even arrives
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
