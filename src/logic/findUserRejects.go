package logic

import (
	"calendar-synch/src/objects"
	"calendar-synch/src/logic/my_datastore"
	"context"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/appengine"
	"log"
	gaeLog "google.golang.org/appengine/log"
	"sort"
	"google.golang.org/appengine/urlfetch"
)

// FIXME this is a major temporary hack
func FindUserRejects(ctx context.Context, cal *calendar.Service) ([]*objects.Event, error) {
	saved, err := my_datastore.QueryEvents(ctx)
	if err != nil {
		return nil, err
	}
	actual, err := cal.Events.List("primary").Do()
	if err != nil {
		return nil, err
	}
	savedSortable := objects.SortableEvents(saved)
	actualSortable := objects.SortableEvents(EventsMap(actual.Items, objects.ConvertGoogleEventToMyEvent))

	if appengine.IsDevAppServer() {
		log.Printf("\nSaved: %v\n", savedSortable)
		log.Printf("\nActual: %v\n", actualSortable)
	} else {
		gaeLog.Debugf(ctx, "\nSaved: %v\n", savedSortable)
		gaeLog.Debugf(ctx, "\nActual: %v\n", actualSortable)
	}

	sort.Sort(savedSortable)  // S, i indices
	sort.Sort(actualSortable) // A, j indices
	lenS := len(savedSortable)
	lenA := len(actualSortable)

	if lenS != lenA {
		whereTo := "https://calendarcron.appspot.com/notify/get"
		client := urlfetch.Client(ctx)
		resp, err := client.Get(whereTo)
		if err != nil {
			gaeLog.Debugf(ctx, "Stop hacking so much: %s", err.Error())
		} else {
			gaeLog.Debugf(ctx, "Still stop hacking so much: %s", resp.Status)
		}
		return nil, err
	}

	for i:=0; i<lenS; i++ {
		//s := savedSortable[i]
		//a := actualSortable[i]
		//
		//s
	}

	return nil, nil
}
