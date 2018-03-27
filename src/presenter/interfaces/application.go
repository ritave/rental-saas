package interfaces

import (
	"rental-saas/src/model"
	"rental-saas/src/model/my_datastore"
	"google.golang.org/api/calendar/v3"
	"time"
)

type DatastoreInterface interface{
	QueryEvents() ([]*model.Event, error)
	DeleteEvent(UUID string) (error)
	SaveEvent(event *model.Event) (error)
	SynchroniseDatastore([]*model.EventModified) (my_datastore.SynchEffect)
}
type CalendarInterface interface{
	AddEvent(event *model.Event) (*model.Event, error)
	DeleteEvent(UUID string) (error)
	QueryEvents() ([]*model.Event, error)
	UpdateEvent(UUID string, event *calendar.Event) (error)
	WatchForChanges(receiver string, expireAfter time.Duration) (error)
}