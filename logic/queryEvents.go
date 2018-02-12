package logic

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"calendar-synch/objects"
)

//QueryEvents returns all the events in datastorage.
func QueryEvents(ctx context.Context) ([]*objects.Event, error) {
	// Print out previous events.
	q := datastore.NewQuery("Event")

	events := make([]*objects.Event, 0)
	_, err := q.GetAll(ctx, &events)
	return events, err
}
