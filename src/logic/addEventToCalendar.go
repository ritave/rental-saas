package logic

import (
	"google.golang.org/api/calendar/v3"
	"log"
	"calendar-synch/src/objects"
	"calendar-synch/src/utils"
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

	//creationDate := utils.TimeToString(time.Now())

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
		return time.Parse(utils.DefaultTimeType, in)
	}

	startT, err := stringToTime(ev.Start)
	if err != nil {
		log.Printf("Start date: %s", err.Error())
		return errors.New("invalid datetime format (accepted is RFC3339: 2006-01-02T15:04:05Z or 2006-01-02T15:04:05+07:00); supplied was: "+ev.Start)
	}

	endT, err := stringToTime(ev.End)
	if err != nil {
		log.Printf("End date: %s", err.Error())
		return errors.New("invalid datetime format (accepted is RFC3339: 2006-01-02T15:04:05Z or 2006-01-02T15:04:05+07:00); supplied was: "+ev.End)
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