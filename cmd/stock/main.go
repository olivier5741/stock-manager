package stock

import (
	. "github.com/olivier5741/stock-manager/skelet"
)

var (
	Chain = []func(cmd Cmd) Cmd{
		Route, Error,
		Get, Error,
		Act, Error,
		Put, Error,
		Publish, Error,
	}
	// On peut même rajouter un ESB :) ou n'importe quoi !!
)

type EndPt struct {
	Db EvtSrcPersister
}
