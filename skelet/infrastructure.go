package skelet

import (
	"fmt"
	"log"
)

type EvtSrcPersister interface {
	GetAll() ([]Ider, error) // add error
	Get(id string) (Ider, error)
	Put(id string, event interface{}) error
}

type Ider interface {
	Id() string
}

type AggAct func(agg interface{}, cmd interface{}) (event interface{}, extEvent interface{}, err error)

type Cmd struct {
	T        Ider
	Agg      Ider
	Act      AggAct
	Event    interface{}
	ExtEvent interface{}
	Route    func(t Ider) (ok bool, a AggAct, p EvtSrcPersister)
	Persist  EvtSrcPersister
	Err      error
}

func ExecuteCommand(cmd Cmd, chain []func(cmd Cmd) Cmd) {
	for _, link := range chain {
		cmd = link(cmd)
	}
}

func Error(cmd Cmd) Cmd {
	if cmd.Err != nil {
		log.Println(cmd.Err)
		panic(cmd.Err)
	}
	return cmd
}

func Get(cmd Cmd) Cmd {
	cmd.Agg, cmd.Err = cmd.Persist.Get(cmd.T.Id())
	log.Println(cmd.Agg)
	return cmd
}

func Act(cmd Cmd) Cmd {
	p := cmd.Agg
	cmd.Event, cmd.ExtEvent, cmd.Err = cmd.Act(p, cmd.T)
	return cmd
}

func Put(cmd Cmd) Cmd {
	cmd.Err = cmd.Persist.Put(cmd.Agg.Id(), cmd.Event)
	//	log.Println(cmd.Event)
	//	log.Println(cmd.Agg)
	return cmd
}

func Route(cmd Cmd) Cmd {
	if cmd.Route == nil {
		return cmd
	}

	ok, act, per := cmd.Route(cmd.T)
	cmd.Act = act
	cmd.Persist = per
	if !ok {
		cmd.Err = fmt.Errorf("Cannot find the route")
	}

	return cmd
}

func Publish(cmd Cmd) Cmd {
	//	log.Print("Event Published : ")
	//	log.Println(cmd.ExtEvent)
	return cmd
}
