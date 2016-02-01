package stock

import (
	sk "github.com/olivier5741/stock-manager/skelet"
)

var (
	Chain = []func(cmd sk.Cmd) sk.Cmd{
		sk.Route, sk.Error,
		sk.Get, sk.Error,
		sk.Act, sk.Error,
		sk.Put, sk.Error,
		sk.Publish, sk.Error,
	}
	// On peut même rajouter un ESB :) ou n'importe quoi !!
)

type EndPt struct {
	Db sk.EvtSrcPersister
}
