package logic

import (
	"google.golang.org/api/calendar/v3"
	"golang.org/x/net/context"
	"time"
	"strconv"
	"sort"
)

type EventModified struct {
	Event *Event
	Action EventModification
}

type EventModification uint

const (
	Deleted EventModification = iota
	ChangedDate EventModification = iota
	ChangedLocaation EventModification = iota
)

func FindChanged(ctx context.Context, cal *calendar.Service) ([]EventModified, error) {
	saved, err := QueryEvents(ctx)
	if err != nil {
		return nil, err
	}
	actual, err := cal.Events.List("primary").Do()
	if err != nil {
		return nil, err
	}
	savedSortable := SortableEvents(saved)
	actualSortable := SortableEvents(EventsMap(actual.Items, ConvertEventToEventLol))

	sort.Sort(SortableEvents(savedSortable))
	sort.Sort(SortableEvents(actualSortable))

	return Compare(savedSortable, actualSortable)
}

func Compare(saved SortableEvents, actual SortableEvents) ([]EventModified, error) {
	panic("implement me")
}

func ConvertEventToEventLol(gEvent *calendar.Event) (myEvent *Event, err error) {
	myEvent = &Event{}

	// user, what if user added someone as attendee or rejected being invited to it?
	if len(gEvent.Attendees) != 1 {
		return myEvent, ConvertingErrorConstructor(UserScrewedUpTheEvent)
	} else {
		myEvent.User = gEvent.Attendees[0].Email
	}

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

func EventsMap(vs []*calendar.Event, f func(event *calendar.Event) (*Event, error)) []*Event {
	vsm := make([]*Event, len(vs))
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
