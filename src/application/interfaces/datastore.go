package interfaces

import (
	"rental-saas/src/model"
	"rental-saas/src/model/my_datastore"
)

// TODO split 'em

type DatastoreInterface interface{
	QueryEvents() ([]*model.Event, error)
	GetEvent(UUID string) (*model.Event, error)
	DeleteEvent(UUID string) (error)
	SaveEvent(event *model.Event) (error)
	Restart() ()
	SynchroniseDatastore([]*model.EventModified) (my_datastore.SynchEffect)
}
