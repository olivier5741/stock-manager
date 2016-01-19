package skelet

import (
	log "github.com/Sirupsen/logrus"
)

type EvtPlay func(acts []interface{}, id string) Ider
type Init func(id string) Ider

type DumEvtRepo struct {
	db      map[string][]interface{}
	InitAgg Init
	PlayEvt EvtPlay
}

func MakeDumEvtRepo(i Init, e EvtPlay) *DumEvtRepo {
	db := map[string][]interface{}{}
	return &DumEvtRepo{db, i, e}
}

func (r DumEvtRepo) GetAll() (aggs []Ider, err error) {
	for key, _ := range r.db {
		out, err := r.Get(key)
		log.Debug("AGG AFTER GET")
		log.Debug(out)
		if err != nil {
			return nil, err
		}
		aggs = append(aggs, out)
	}
	return
}

func (r DumEvtRepo) Get(id string) (Ider, error) {
	events, ok := r.db[id]
	if !ok {
		return r.InitAgg(id), nil
	}
	return r.PlayEvt(events, id), nil
}

func (r *DumEvtRepo) Put(id string, event interface{}) (err error) {
	log.Debug("EVENTS TO PUT")
	log.Debug(event)
	r.db[id] = append(r.db[id], event)
	return
}
