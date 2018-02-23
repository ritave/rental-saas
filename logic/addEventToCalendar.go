package logic

import (
	"google.golang.org/api/calendar/v3"
	"log"
	"calendar-synch/objects"
	"calendar-synch/helpers"
	"errors"
	"regexp"
	"time"
)


func AddEventToCalendar(cal *calendar.Service, ev objects.Event) (*objects.Event, error){
	newEvent := &calendar.Event{
		Summary:     ev.Summary,
		Location:    ev.Location,
		Description: "Cleaning service ordered on %s. Feel free to move this event in your calendar to change the date!",
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

	//creationDate := helpers.TimeToString(time.Now())

	evResp, err := cal.Events.Insert("primary", newEvent).Do()
	if err != nil {
		log.Println("Adding event failed")
		log.Printf("Error: %s", err.Error())
		log.Printf("Event: %v", ev)
	} else {
		log.Printf("Link: %s", evResp.HtmlLink)

		// Creation date will be my "primary-key"
		eventPrimaryKey := evResp.Created
		ev.CreationDate = eventPrimaryKey
		return &ev, nil
	}

	return nil, err
}

var emailParser = regexp.MustCompile(`(^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$)`)

func EvenMoreChecksForTheEvent(ev objects.Event) (error) {

	// time checks
	var stringToTime = func(in string) (time.Time, error) {
		return time.Parse(helpers.DefaultTimeType, in)
	}

	startT, err := stringToTime(ev.Start)
	if err != nil {
		return errors.New("invalid error format")
	}

	endT, err := stringToTime(ev.End)
	if err != nil {
		return errors.New("invalid error format")
	}

	if startT.Before(time.Now()) {
		return errors.New("event start cannot be set in the past")
	}

	if endT.Before(startT) {
		return errors.New("event end cannot be before the start")
	}

	if endT.Before(time.Now()) {
		return errors.New("event end cannot be set in the past")
	}

	if ev.Location == "" {
		return errors.New("location cannot be empty")
	}

	if !emailParser.Match([]byte(ev.User)) {
		return errors.New("invalid email address")
	}

	return nil
}