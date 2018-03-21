package my_datastore

import (
	"rental-saas/src/model"
)

type Datastore struct {

}

func (ds *Datastore) QueryEvents() ([]*model.Event, error) {
	panic("implement me")
}

func (ds *Datastore) DeleteEvent(UUID string) (error) {
	panic("implement me")
}

func (ds *Datastore) SaveEvent(event *model.Event) (error) {
	panic("implement me")
}

func (ds *Datastore) SynchroniseDatastore() {
	panic("implement me")
}

