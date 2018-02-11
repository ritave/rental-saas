package logic

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

type Event struct {
	Summary  string
	User     string
	Start    string
	End      string
	Location string
}

// TODO ancestors

func SaveEventInDatastore(ctx context.Context, ev Event) error {
	k := datastore.NewIncompleteKey(ctx,"Event", nil)

	_, err := datastore.Put(ctx, k, &ev)
	return err
}

//QueryEvents returns all the events in datastorage.
func QueryEvents(ctx context.Context) ([]*Event, error) {
	// Print out previous events.
	q := datastore.NewQuery("Event")

	events := make([]*Event, 0)
	_, err := q.GetAll(ctx, &events)
	return events, err
}
