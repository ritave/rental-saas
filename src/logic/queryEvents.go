package logic

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"calendar-synch/src/objects"
	"time"
	"log"
	"calendar-synch/src/utils"
)

const EventKeyKind = "Event"

//QueryEvents returns all the events in datastorage.
func QueryEvents(ctx context.Context) ([]*objects.Event, error) {
	// Print out previous events.
	q := datastore.NewQuery(EventKeyKind)

	events := make([]*objects.Event, 0)
	_, err := q.GetAll(ctx, &events)
	return events, err
}

func QueryEventsFiltered(ctx context.Context) ([]*objects.Event, error) {
	now := utils.TimeToMilliseconds(time.Now())
	log.Printf("I want events that happenend before: %d", now)
	q := datastore.NewQuery(EventKeyKind).Filter("Timestamp <", now)

	events := make([]*objects.Event, 0)
	_, err := q.GetAll(ctx, &events)
	return events, err
}
