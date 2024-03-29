package core

import (
	"rental-saas/src/utils/config"
	"rental-saas/src/utils"
	"rental-saas/src/application/interfaces"
	"rental-saas/src/model/my_datastore"
	"rental-saas/src/model/my_calendar"
	"rental-saas/src/api_integration"
	"net/http"
)

type Application struct {
	Datastore interfaces.DatastoreInterface
	Calendar  interfaces.CalendarInterface
	Config    config.C

	// hidden for less clutter
	Utils Utils
}
type Utils struct {
	Ticker *utils.Ticker
	Pozamiatane api_integration.Provider
}

func New(c config.C) *Application {
	return &Application{
		Datastore: my_datastore.New(c),
		Calendar: my_calendar.New(c),
		Config: c,

		Utils: Utils{
			Pozamiatane: api_integration.Provider{
				Client: http.DefaultClient,
				Server: c.Pozamiatane.Address,
			},
		},
	}
}
