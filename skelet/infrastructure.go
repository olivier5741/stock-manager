package skelet

type UniquePersister interface {
	Get() (agg interface{})
	Put(agg interface{})
}

type EvtSrcPersister interface {
	GetAll() (aggs []interface{})
	Get(id string) (agg interface{})
	Put(id string, event interface{})
}
