package my_calendar

import (
	"rental-saas/src/model"
)

type Calendar struct {

}

func (c *Calendar) AddEvent(event *model.Event) (*model.Event, error) {
	panic("implement me")
}

func (c *Calendar) DeleteEvent(UUID string) (error) {
	panic("implement me")
}


