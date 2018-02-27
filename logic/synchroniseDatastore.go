package logic

import (
	"calendar-synch/objects"
	"context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"errors"
)

func SynchroniseDatastore(ctx context.Context, diff []*objects.EventModified) {
	for _, event := range diff {
		for mod := range event.Modifications {
			if mod == objects.Deleted {
				// event has been deleted from the calendar, it has to be deleted from the datastore

				err := deleteEvent(ctx, event.Event)
				if err != nil {
					log.Criticalf(ctx,"SYNCHRONISE | Event deleting FAILED: %s", err.Error())
				}
				// we do not have to do anything more, carry on with the outermost for loop
				break
			} else if mod == objects.Added {
				// event added, add to datastore

				err := SaveEventInDatastore(ctx, event.Event)
				if err != nil {
					log.Criticalf(ctx, "SYNCHRONISE | Event adding FAILED: %s", err.Error())
				}
				// TODO added and modified possible?
			} else if mod == objects.ModifiedTime || mod == objects.ModifiedLocation {
				// modified -> delete and add again

				// deleting uses only the UUID and I think I'm copying it from old_modified to new_modified in findChanged()
				err := deleteEvent(ctx, event.Event)
				if err != nil {
					log.Criticalf(ctx,"SYNCHRONISE | Event deleting FAILED: %s", err.Error())
				}
				err = SaveEventInDatastore(ctx, event.Event)
				if err != nil {
					log.Criticalf(ctx, "SYNCHRONISE | Event adding FAILED: %s", err.Error())
				}
			}

		}
	}
}

func deleteEvent(ctx context.Context, event *objects.Event) (error) {
	if event.UUID == "" {
		return errors.New("event UUID was empty when it was needed most")
	}

	eventKey := datastore.NewKey(ctx, EventKeyKind, event.UUID, 0, nil)
	return datastore.Delete(ctx, eventKey)
}