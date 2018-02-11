package logic

import (
	"google.golang.org/api/calendar/v3"
	"time"
	"log"
)

type Event struct {
	Summary  string
	User     string
	Start    string
	End      string
	Location string
}

func AddEventToCalendar(cal *calendar.Service, eventRequest Event) {
	newEvent := &calendar.Event{
		Summary:     eventRequest.Summary,
		Location:    eventRequest.Location,
		Description: "Cleaning service ordered on %s. Feel free to move this event in your calendar to change the date!",
		Start: &calendar.EventDateTime{
			DateTime: time.Now().Format(time.RFC3339),
		},
		End: &calendar.EventDateTime{
			DateTime: time.Now().Add(time.Hour).Format(time.RFC3339),
		},
		Attendees: []*calendar.EventAttendee{
			&calendar.EventAttendee{Email: eventRequest.User},
		},
		GuestsCanModify: true,
	}

	ev, err := cal.Events.Insert("primary", newEvent).Do()
	if err != nil {
		log.Println("Adding event failed")
		log.Printf("Error: %s", err.Error())
	} else {
		log.Printf("Link: %s", ev.HtmlLink)
	}
}
