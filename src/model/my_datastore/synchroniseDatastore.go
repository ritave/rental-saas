package my_datastore

import (
	"rental-saas/src/model"
	"context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"errors"
)

type SynchEffect struct {
	Errors int
	Deleted int
	Added int
	Modified int
}

func SynchroniseDatastore(ctx context.Context, diff []*model.EventModified) (SynchEffect) {
	var result SynchEffect
	
	for _, event := range diff {
		for mod := range event.Modifications {
			if mod == model.Deleted {
				// event has been deleted from the calendar, it has to be deleted from the datastore

				err := deleteEvent(ctx, event.Event)
				if err != nil {
					log.Debugf(ctx,"SYNCHRONISE | Event deleting FAILED: %s", err.Error())
					result.Errors ++
					break
				}
				// we do not have to do anything more, carry on with the outermost for loop
				result.Deleted ++
				break
			} else if mod == model.Added {
				// event added, add to datastore

				err := SaveEventInDatastore(ctx, event.Event)
				if err != nil {
					log.Debugf(ctx, "SYNCHRONISE | Event adding FAILED: %s", err.Error())
					result.Errors ++
					break
				}

				result.Added ++

				// TODO added and modified possible?
			} else if mod == model.ModifiedTime || mod == model.ModifiedLocation {
				// modified -> delete and add again

				// deleting uses only the UUID and I think I'm copying it from old_modified to new_modified in findChanged()
				err := deleteEvent(ctx, event.Event)
				if err != nil {
					log.Debugf(ctx,"SYNCHRONISE | Event deleting FAILED: %s", err.Error())
					result.Errors ++
					break
				}
				err = SaveEventInDatastore(ctx, event.Event)
				if err != nil {
					log.Debugf(ctx, "SYNCHRONISE | Event adding FAILED: %s", err.Error())
					result.Errors ++
					break
				}

				result.Modified ++
			}

		}
	}

	return result
}

func deleteEvent(ctx context.Context, event *model.Event) (error) {
	if event.UUID == "" {
		return errors.New("event UUID was empty when it was needed most")
	}

	eventKey := datastore.NewKey(ctx, EventKeyKind, event.UUID, 0, nil)
	return datastore.Delete(ctx, eventKey)

}