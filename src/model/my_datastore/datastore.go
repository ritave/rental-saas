package my_datastore

import (
	"rental-saas/src/model"
)

// TODO
// major TODO:
// get some kind of a database: MySQL, SQLite?
// worst-case scenario would be Mongo

type Datastore struct {
	// not really a persistent database, lol
	ds map[string]*model.Event
}

func (ds *Datastore) SynchroniseDatastore([]*model.EventModified) (SynchEffect) {
	panic("implement me")
}

func (ds *Datastore) QueryEvents() ([]*model.Event, error) {
	result := make([]*model.Event, len(ds.ds))
	i := 0
	for _, v := range ds.ds {
		result[i] = v
		i ++
	}
	return result, nil
}

func (ds *Datastore) DeleteEvent(UUID string) (error) {
	ds.ds[UUID] = nil
	return nil
}

func (ds *Datastore) SaveEvent(event *model.Event) (error) {
	ds.ds[event.UUID] = event
	return nil
}

