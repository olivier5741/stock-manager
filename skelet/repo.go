package skelet

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
		if err != nil {
			return nil, err
		}
		aggs = append(aggs, out)
	}
	return
}

func (r DumEvtRepo) GetAllEvents(id string) (events []interface{}) {
	events, ok := r.db[id]
	if !ok {
		events = make([]interface{}, 0) // not sure I need this
	}
	return
}

func (r DumEvtRepo) Get(id string) (Ider, error) {
	events, ok := r.db[id]
	if !ok {
		return r.InitAgg(id), nil
	}

	acts := make([]interface{}, 0)
	for _, e := range events {
		// TODO : something more intelligent than this :)
		acts = append(acts, e.(Event).Act)
	}

	return r.PlayEvt(acts, id), nil
}

func (r *DumEvtRepo) Put(id string, event interface{}) (err error) {
	r.db[id] = append(r.db[id], event)
	return
}
