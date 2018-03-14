package objects

import (
	"calendar-synch/src/utils"
	"google.golang.org/api/calendar/v3"
	"strconv"
	"strings"
	"fmt"
)

func ConvertGoogleToMine(gEvent *calendar.Event) (myEvent *Event) {
	myEvent = &Event{}

	// user, what if user added someone as attendee or rejected being invited to it?
	probablyManyEmails := make([]string, len(gEvent.Attendees))
	for ind, at := range gEvent.Attendees {
		probablyManyEmails[ind] = at.Email
	}
	myEvent.User = strings.Join(probablyManyEmails, ";")

	// creation date
	creation := utils.StringToTime(gEvent.Created)
	myEvent.Timestamp = utils.TimeToMilliseconds(creation)
	myEvent.CreationDate = gEvent.Created

	// date
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

	// testFields
	var addField = func(key string, val interface{}) (string) {
		return fmt.Sprintf("%s %v; ", key, val)
	}
	var testString string
	testString += addField("AttendeeOmitted", gEvent.AttendeesOmitted)
	testString += addField("Status", gEvent.Status)
	testString += addField("GuestsCanInviteOthers", gEvent.GuestsCanInviteOthers)
	testString += addField("GuestsCanModify", gEvent.GuestsCanModify)
	if len(gEvent.Attendees) > 0 {
		testString += addField("Attendees", len(gEvent.Attendees))
		testString += addField("ResponseStatus", gEvent.Attendees[0].ResponseStatus)
		testString += addField("Self", gEvent.Attendees[0].Self)
	}

	return myEvent
}

func ConvertGoogleToMineSlice(vs []*calendar.Event) []*Event {
	vsm := make([]*Event, len(vs))
	for i, v := range vs {
		vsm[i] = ConvertGoogleToMine(v)
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
