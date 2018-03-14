package logic

import (
	"google.golang.org/api/calendar/v3"
	"golang.org/x/net/context"
	"sort"
	"calendar-synch/src/objects"
	"google.golang.org/appengine"
	"log"
	"calendar-synch/src/logic/my_datastore"
	gaeLog "google.golang.org/appengine/log"
)


func FindChanged(ctx context.Context, cal *calendar.Service) ([]*objects.EventModified, error) {
	saved, err := my_datastore.QueryEvents(ctx)
	if err != nil {
		return nil, err
	}
	actual, err := cal.Events.List("primary").Do()
	if err != nil {
		return nil, err
	}
	savedSortable := objects.SortableEvents(saved)
	actualSortable := objects.SortableEvents(objects.ConvertGoogleToMineSlice(actual.Items))

	if appengine.IsDevAppServer() {
		log.Printf("\nSaved: %v\n", savedSortable)
		log.Printf("\nActual: %v\n", actualSortable)
	} else {
		gaeLog.Debugf(ctx, "\nSaved: %v\n", savedSortable)
		gaeLog.Debugf(ctx, "\nActual: %v\n", actualSortable)
	}

	return CompareSortable(savedSortable, actualSortable, ctx)
}

func CompareSortable(saved objects.SortableEvents, actual objects.SortableEvents, ctx context.Context) ([]*objects.EventModified, error) {
	// sort by creation date
	sort.Sort(objects.SortableEvents(saved)) // S, i indices
	sort.Sort(objects.SortableEvents(actual)) // A, j indices

	lenS := len(saved)
	lenA := len(actual)

	changes := make([]*objects.EventModified, 0)

	var i, j int
	for i < lenS && j < lenA {
		s := saved[i]
		a := actual[j]

		/*
		  ---time--->
		S: [ ]   [ ][ ]
		A:    [ ][ ][.]
		    d  a     m
		Where: (d)eleted, (a)dded, (m)odified

		 */

		if s.Equal(a) {
			// present in saved and actual => MAYBE actual has been modified?
			if !s.IsTheSame(a) {
				// they are not the same, find out what has changed

				// create and append the object now, we will flag it later (maybe even multiple times)
				d := objects.NewModified(a)

				// CARRY OVER THE UUID SO WE CAN REFERENCE IT LATER WHEN SYNCHRONISING
				d.Event.UUID = s.UUID


				modifications := 0

				if s.Location != a.Location {
					d.Flag(objects.ModifiedLocation)
					modifications ++
				}

				if s.Start != a.Start || s.End != a.End {
					d.Flag(objects.ModifiedTime)
					modifications ++
				}

				if a.User != s.User {
					if a.User == "" {
						d.Flag(objects.UserRejected)
					} else {
						d.Flag(objects.SomethingWonkyHappened)
						gaeLog.Debugf(ctx, "WONKY lol: Actual %#v; Saved %#v", a, s)
					}
				}

				// TODO
				// to be added more later...

				// so the keyword here is "MAYBE"
				if modifications > 0 {
					changes = append(changes, d)
				}
				// otherwise event was not modified at all
			}

			// they were the same event => move both indices forward
			i++
			j++

		} else if a.Less(s) { // this should happen more often than the last clause
			// "a" is missing in saved => "a" has been added
			changes = append(changes, objects.NewModified(a).Flag(objects.Added))

			// a is behind s => next a
			j++
		} else {
			// "s" is missing in actual => "s" has been deleted
			changes = append(changes, objects.NewModified(s).Flag(objects.Deleted))

			// s is behind a => next s
			i++
		}
	}

	// what should remain now is ONLY elements of S XOR ONLY elements of A

	// exhaust elements of S
	for ; i < lenS; i++ {
		// => deleted
		s := saved[i]
		changes = append(changes, objects.NewModified(s).Flag(objects.Deleted))
	}

	// exhaust elements of A
	for ; j < lenA; j++ {
		// => added
		a := actual[j]
		changes = append(changes, objects.NewModified(a).Flag(objects.Added))
	}

	return changes, nil
}
