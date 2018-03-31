package my_datastore

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"rental-saas/src/model"
	"time"
		"rental-saas/src/utils"
	"github.com/sirupsen/logrus"
)

const EventKeyKind = "Event"

//QueryEvents returns all the events in datastorage.
func QueryEvents(ctx context.Context) ([]*model.Event, error) {
	// Print out previous events.
	q := datastore.NewQuery(EventKeyKind)

	events := make([]*model.Event, 0)
	_, err := q.GetAll(ctx, &events)
	return events, err
}

func QueryEventsFiltered(ctx context.Context) ([]*model.Event, error) {
	now := utils.TimeToMilliseconds(time.Now())
	logrus.Printf("I want events that happenend before: %d", now)
	q := datastore.NewQuery(EventKeyKind).Filter("Timestamp <", now)

	events := make([]*model.Event, 0)
	_, err := q.GetAll(ctx, &events)
	return events, err
}
