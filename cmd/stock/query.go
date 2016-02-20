package stock

import (
	"github.com/olivier5741/stock-manager/item/items"
	"github.com/olivier5741/stock-manager/skelet"
	"github.com/olivier5741/stock-manager/stock"
	"fmt"
)

func (endPt EndPt) ProdValEvol(id string) (data map[string]items.Items) {
	// TODO : should be generated when even arrives
	acts := endPt.Db.GetAllEvents(id)

	// REVERSE (don't know why I have to do this)
	for i, j := 0, len(acts)-1; i < j; i, j = i+1, j-1 {
        acts[i], acts[j] = acts[j], acts[i]
    }

	data = make(map[string]items.Items, len(acts))
	state := make(items.Items, 0)

	i := 0
loop:
	for _, a := range acts {
		fmt.Println(a)
		switch act := a.(skelet.Event).Act.(type) {
		case stock.In:
			state = items.Add(state, act.Items)
		case stock.Out:
			state = items.Sub(state, act.Items)
		case stock.Inventory:
			state = act.Items.Copy()
		default:
			continue loop
		}
		data[a.(skelet.Event).Date] = state.Copy()
		i++
	}

	return

}
