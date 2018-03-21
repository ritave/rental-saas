package wrapper

import (
	"rental-saas/src/utils/config"
	"rental-saas/src/utils"
	"rental-saas/src/model"
	"rental-saas/src/model/my_datastore"
	"google.golang.org/api/calendar/v3"
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
}

type Application struct {
	Datastore DatastoreInterface
	Calendar  CalendarInterface
	Config    config.C

	// hidden for less clutter
	Utils *Utils
}
type Utils struct {
	Ticker *utils.Ticker
}

func New(c config.C) *Application {
	return &Application{
	}
}

