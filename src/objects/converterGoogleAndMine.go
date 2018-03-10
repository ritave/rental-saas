package objects

import (
	"calendar-synch/src/utils"
	"time"
	"google.golang.org/api/calendar/v3"
	"strconv"
	"fmt"
)

func ConvertGoogleToMine(gEvent *calendar.Event) (myEvent *Event, err error) {
	myEvent = &Event{}

	// TODO
	// user, what if user added someone as attendee or rejected being invited to it?
	if len(gEvent.Attendees) != 1 {
		return myEvent, ConvertingErrorConstructor(UserScrewedUpTheEvent)
	} else {
		myEvent.User = gEvent.Attendees[0].Email
	}

	// creation date
	creation := utils.StringToTime(gEvent.Created)
	myEvent.Timestamp = utils.TimeToMilliseconds(creation)
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

	// uuid
	myEvent.UUID = gEvent.Id

	// timestamp
	myEvent.Timestamp = utils.TimeToMilliseconds(utils.StringToTime(gEvent.Created))

	// creationDate
	myEvent.CreationDate = gEvent.Created

	return myEvent, err
}

func ConvertGoogleToMineSlice(vs []*calendar.Event) []*Event {
	vsm := make([]*Event, len(vs))
	for i, v := range vs {
		vsm[i], _ = ConvertGoogleToMine(v) // this shouldn't return an error, 'cause Google is the smarter one here
	}
	return vsm
}

func ConvertMineToGoogle(ev Event) (*calendar.Event) {


	newEvent := &calendar.Event{
		Summary:     ev.Summary,
		Location:    ev.Location,
		Description: fmt.Sprintf("Cleaning service ordered on %s. Feel free to move this event in your calendar to change the date!", time.Now().Format(time.RFC822)),
		Start: &calendar.EventDateTime{
			DateTime: ev.Start,
		},
		End: &calendar.EventDateTime{
			DateTime: ev.End,
		},
		Attendees: []*calendar.EventAttendee{
			{Email: ev.User},
		},
		GuestsCanModify: true, // that's what allows for changing the date of the event... but also all the other fields
	}

	return newEvent
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
