package logic

import (
	"google.golang.org/api/calendar/v3"
	"time"
	"log"
	"calendar-synch/objects"
)


func AddEventToCalendar(cal *calendar.Service, ev objects.Event) (*objects.Event, error){
	newEvent := &calendar.Event{
		Summary:     ev.Summary,
		Location:    ev.Location,
		Description: "Cleaning service ordered on %s. Feel free to move this event in your calendar to change the date!",
		Start: &calendar.EventDateTime{
			DateTime: time.Now().Format(time.RFC3339), // FIXME temporary
		},
		End: &calendar.EventDateTime{
			DateTime: time.Now().Add(time.Hour).Format(time.RFC3339), // FIXME temporary
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
