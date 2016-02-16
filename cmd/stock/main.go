package stock

import (
	"github.com/olivier5741/stock-manager/skelet"
)

var (
	Chain = []func(cmd skelet.Cmd) skelet.Cmd{
		skelet.Route, skelet.Error,
		skelet.Get, skelet.Error,
		skelet.Act, skelet.Error,
		skelet.Put, skelet.Error,
		skelet.Publish, skelet.Error,
	}
	// On peut même rajouter un ESB :) ou n'importe quoi !!
)

type EndPt struct {
	Db skelet.EvtSrcPersister
}
