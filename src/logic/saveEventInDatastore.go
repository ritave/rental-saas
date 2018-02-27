package logic

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"calendar-synch/src/objects"
	"github.com/satori/go.uuid"
)

func SaveEventInDatastore(ctx context.Context, ev *objects.Event) error {
	// this will come in handy when we want to delete the event from datastore
	ev.UUID = uuid.Must(uuid.NewV4()).String()

	//k := datastore.NewIncompleteKey(ctx,"Event", nil)
	k := datastore.NewKey(ctx, EventKeyKind, ev.UUID, 0, nil)

	_, err := datastore.Put(ctx, k, ev)
	return err
}
