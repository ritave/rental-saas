package wrapper

import (
	"rental-saas/src/utils/config"
	"rental-saas/src/utils"
	"rental-saas/src/presenter/interfaces"
	"rental-saas/src/model/my_datastore"
	"rental-saas/src/model/my_calendar"
)

type Application struct {
	Datastore interfaces.DatastoreInterface
	Calendar  interfaces.CalendarInterface
	Config    config.C

	// hidden for less clutter
	Utils *Utils
}
type Utils struct {
	Ticker *utils.Ticker
}

func New(c config.C) *Application {
	return &Application{
		Datastore: my_datastore.New(c),
		Calendar: my_calendar.New(c),
	}
}
