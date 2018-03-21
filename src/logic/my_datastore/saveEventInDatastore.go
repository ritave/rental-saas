package my_datastore

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"rental-saas/src/objects"
)

func SaveEventInDatastore(ctx context.Context, ev *objects.Event) error {
	//k := datastore.NewIncompleteKey(ctx,"Event", nil)
	k := datastore.NewKey(ctx, EventKeyKind, ev.UUID, 0, nil)

	_, err := datastore.Put(ctx, k, ev)
	return err
}
