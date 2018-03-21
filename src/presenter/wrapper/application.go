package wrapper

import (
	"rental-saas/src/utils/config"
	"rental-saas/src/utils"
	"rental-saas/src/model"
)

type DBInterface interface{
	QueryEvents()
	DeleteEvent(event *model.Event)
	SaveEvent(event *model.Event) (error)
	SynchroniseDatastore()
}
type CalendarInterface interface{
	AddEvent(event *model.Event) (*model.Event, error)
	DeleteEvent(event *model.Event) (error)
}

type Application struct {
	DB DBInterface
	Calendar CalendarInterface
	Config config.C

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

