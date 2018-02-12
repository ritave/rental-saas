package logic

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"calendar-synch/objects"
)

func SaveEventInDatastore(ctx context.Context, ev objects.Event) error {
	k := datastore.NewIncompleteKey(ctx,"Event", nil)

	_, err := datastore.Put(ctx, k, &ev)
	return err
}
