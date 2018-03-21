package wrapper

import (
	"rental-saas/src/utils/config"
	"rental-saas/src/utils"
)

type DBInterface interface{}
type CalendarInterface interface{}

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

