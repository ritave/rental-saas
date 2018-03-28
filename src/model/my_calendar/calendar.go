package my_calendar

import (
	"rental-saas/src/model"
	"time"
	"google.golang.org/api/calendar/v3"
	"rental-saas/src/presenter"
	"rental-saas/src/utils/config"
	"rental-saas/src/calendar_wrap"
	"context"
)

type Calendar struct {
	Service *calendar.Service
}

func New(c config.C) (*Calendar) {
	service := calendar_wrap.NewFlex(context.Background())
	return &Calendar{
		Service: service,
	}
}

func (c *Calendar) QueryEvents() ([]*model.Event, error) {
	gEvents, err := c.Service.Events.List("primary").Do()
	if err != nil {
		return nil, err
	}
	return model.ConvertGoogleToMineSlice(gEvents.Items), nil
}

func (c *Calendar) UpdateEvent(UUID string, event *calendar.Event) (error) {
	_, err := c.Service.Events.Update("primary", UUID, event).Do()
	return err
}

func (c *Calendar) WatchForChanges(receiver string, expireAfter time.Duration) (error) {
	err, _ := presenter.WatchForChanges(c.Service, receiver, expireAfter)
	return err
}

func (c *Calendar) AddEvent(event *model.Event) (*model.Event, error) {
	return AddEvent(c.Service, event)
}

func (c *Calendar) DeleteEvent(UUID string) (error) {
	return DeleteEvent(c.Service, UUID)
}

