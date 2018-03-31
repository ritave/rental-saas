package interfaces

import (
	"rental-saas/src/model"
	"time"
	"google.golang.org/api/calendar/v3"
)

type CalendarInterface interface{
	AddEvent(event *model.Event) (*model.Event, error)
	DeleteEvent(UUID string) (error)
	QueryEvents() ([]*model.Event, error)
	UpdateEvent(UUID string, event *calendar.Event) (error)
	WatchForChanges(receiver string, expireAfter time.Duration) (error)
}
