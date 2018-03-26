package presenter

import (
	"sort"
	"rental-saas/src/model"
	"google.golang.org/appengine"
	"log"
	"rental-saas/src/presenter/interfaces"
)


func FindChanged(ds interfaces.DatastoreInterface, cal interfaces.CalendarInterface) ([]*model.EventModified, error) {
	saved, err := ds.QueryEvents()
	if err != nil {
		return nil, err
	}

	actual, err := cal.QueryEvents()
	if err != nil {
		return nil, err
	}

	savedSortable := model.SortableEvents(saved)
	actualSortable := model.SortableEvents(actual)

	if appengine.IsDevAppServer() {
		log.Printf("\nSaved: %v\n", savedSortable)
		log.Printf("\nActual: %v\n", actualSortable)
	}

	return CompareSortable(savedSortable, actualSortable)
}

func CompareSortable(saved model.SortableEvents, actual model.SortableEvents) ([]*model.EventModified, error) {
	// sort by creation date
	sort.Sort(model.SortableEvents(saved))  // S, i indices
	sort.Sort(model.SortableEvents(actual)) // A, j indices

	lenS := len(saved)
	lenA := len(actual)

	changes := make([]*model.EventModified, 0)

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
				d := model.NewModified(a)

				// CARRY OVER THE UUID SO WE CAN REFERENCE IT LATER WHEN SYNCHRONISING
				d.Event.UUID = s.UUID


				modifications := 0

				if s.Location != a.Location {
					d.Flag(model.ModifiedLocation)
					modifications ++
				}

				if s.Start != a.Start || s.End != a.End {
					d.Flag(model.ModifiedTime)
					modifications ++
				}

				if a.User != s.User {
					if a.User == "" {
						d.Flag(model.UserRejected)
					} else {
						d.Flag(model.SomethingWonkyHappened)
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
			changes = append(changes, model.NewModified(a).Flag(model.Added))

			// a is behind s => next a
			j++
		} else {
			// "s" is missing in actual => "s" has been deleted
			changes = append(changes, model.NewModified(s).Flag(model.Deleted))

			// s is behind a => next s
			i++
		}
	}

	// what should remain now is ONLY elements of S XOR ONLY elements of A

	// exhaust elements of S
	for ; i < lenS; i++ {
		// => deleted
		s := saved[i]
		changes = append(changes, model.NewModified(s).Flag(model.Deleted))
	}

	// exhaust elements of A
	for ; j < lenA; j++ {
		// => added
		a := actual[j]
		changes = append(changes, model.NewModified(a).Flag(model.Added))
	}

	return changes, nil
}
