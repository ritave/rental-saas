package logic

import (
	"google.golang.org/api/calendar/v3"
	"golang.org/x/net/context"
	"time"
	"strconv"
	"sort"
	"calendar-synch/objects"
)


func FindChanged(ctx context.Context, cal *calendar.Service) ([]*objects.EventModified, error) {
	saved, err := QueryEvents(ctx)
	if err != nil {
		return nil, err
	}
	actual, err := cal.Events.List("primary").Do()
	if err != nil {
		return nil, err
	}
	savedSortable := objects.SortableEvents(saved)
	actualSortable := objects.SortableEvents(EventsMap(actual.Items, ConvertEventToEventLol))

	return CompareSortable(savedSortable, actualSortable)
}

func CompareSortable(saved objects.SortableEvents, actual objects.SortableEvents) ([]*objects.EventModified, error) {
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
		S: [ ]  [ ][ ]
		A:   [ ][ ][.]
		    d a     m
		Where: (d)eleted, (a)dded, (m)odified

		 */

		if s.Equal(a) {
			// present in saved and actual => maybe actual has been modified?
			if !s.IsTheSame(a) {
				// they are not the same, find out what has changed

				// create and append the object now, we will flag it later (maybe even multiple times)
				d := objects.NewModified(a)
				changes = append(changes, d)

				if s.Location != a.Location {
					d.Flag(objects.ModifiedLocation)
				}

				if s.Start != a.Start || s.End != a.End {
					d.Flag(objects.ModifiedTime)
				}

				// to be added more later...
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
		// this exhausts all the cases
	}


	return changes, nil
}

func ConvertEventToEventLol(gEvent *calendar.Event) (myEvent *objects.Event, err error) {
	myEvent = &objects.Event{}

	// user, what if user added someone as attendee or rejected being invited to it?
	if len(gEvent.Attendees) != 1 {
		return myEvent, ConvertingErrorConstructor(UserScrewedUpTheEvent)
	} else {
		myEvent.User = gEvent.Attendees[0].Email
	}

	// creation date
	myEvent.CreationDate = gEvent.Created

	// date
	dtStart, err := time.Parse(time.RFC3339, gEvent.Start.DateTime)
	dtEnd, err := time.Parse(time.RFC3339, gEvent.End.DateTime)
	now := time.Now()
	if dtStart.Before(now) || dtEnd.Before(now) {
		// not fatal I suppose
		err = ConvertingErrorConstructor(DateHasPassed)
	}
	myEvent.Start = gEvent.Start.DateTime
	myEvent.End = gEvent.End.DateTime

	// location
	myEvent.Location = gEvent.Location

	// summary
	myEvent.Summary = gEvent.Summary

	return myEvent, err
}

func EventsMap(vs []*calendar.Event, f func(event *calendar.Event) (*objects.Event, error)) []*objects.Event {
	vsm := make([]*objects.Event, len(vs))
	for i, v := range vs {
		vsm[i], _ = f(v) // LOL xd FIXME eventually
	}
	return vsm
}

const (
	UserScrewedUpTheEvent ConvertingErrorType = iota
	DateHasPassed         ConvertingErrorType = iota
)

func ConvertingErrorConstructor(errorType ConvertingErrorType) (error) {
	return &ConvertingError{
		msg: strconv.Itoa(int(errorType)),
		tp:  errorType,
	}
}

type ConvertingError struct {
	msg string
	tp  ConvertingErrorType
}

type ConvertingErrorType int

func (ce *ConvertingError) Error() (string) {
	return ce.msg
}
