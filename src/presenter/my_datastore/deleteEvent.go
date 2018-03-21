package my_datastore

import (
	"context"
	"google.golang.org/appengine/datastore"
)

func DeleteEvent(ctx context.Context, uuid string) error {
	k := datastore.NewKey(ctx, EventKeyKind, uuid, 0, nil)

	err := datastore.Delete(ctx, k)
	return err
}
