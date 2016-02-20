package skelet

import (
	"fmt"
	// TODO : replace with Campion [golang at FOSDEM] logger
	log "github.com/Sirupsen/logrus"
	//"log"
)

type EvtSrcPersister interface {
	GetAll() ([]Ider, error)
	GetAllEvents(id string) (events []interface{}) // add error later
	Get(id string) (Ider, error)
	Put(id string, event interface{}) error
}

type Ider interface {
	ID() string
}

type Event struct {
	Date string
	Act  interface{}
}

type AggAct func(agg interface{}, cmd interface{}) (event Event, extEvent interface{}, err error)

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
		log.WithFields(log.Fields{
			"err": cmd.Err,
			"cmd": cmd,
		}).Error("Error in command handling")
		log.Println(cmd.Err)
	}
	return cmd
}

func Get(cmd Cmd) Cmd {
	cmd.Agg, cmd.Err = cmd.Persist.Get(cmd.T.ID())
	//log.Debug(cmd.Agg)
	return cmd
}

func Act(cmd Cmd) Cmd {
	p := cmd.Agg
	cmd.Event, cmd.ExtEvent, cmd.Err = cmd.Act(p, cmd.T)
	return cmd
}

func Put(cmd Cmd) Cmd {
	cmd.Err = cmd.Persist.Put(cmd.Agg.ID(), cmd.Event)
	//log.Debug(cmd.Agg)
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
