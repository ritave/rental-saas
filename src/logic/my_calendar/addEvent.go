package my_calendar

import (
	"google.golang.org/api/calendar/v3"
	"log"
	"calendar-synch/src/objects"
	"calendar-synch/src/utils"
	"time"
	"fmt"
	"context"
	gaeLog "google.golang.org/appengine/log"
)


func AddEvent(ctx context.Context, cal *calendar.Service, ev objects.Event) (*objects.Event, error){
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

	//creationDate := utils.TimeToString(time.Now())

	evResp, err := cal.Events.Insert("primary", newEvent).Do()
	if err != nil {
		gaeLog.Debugf(ctx, "Adding event failed %#v %s", ev, err.Error())
	} else {
		log.Printf("Link: %s", evResp.HtmlLink)

		// Creation date will be my "primary-key"
		eventCreationTime, err := utils.VerifyStringToTime(evResp.Created)
		if err != nil {
			gaeLog.Criticalf(ctx, "Google passed to us string that is not of valid format! IMPOSSIBRU!")
			eventCreationTime = time.Now()
		}
		eventOrderingKey := utils.TimeToMilliseconds(eventCreationTime)
		ev.Timestamp = eventOrderingKey
		ev.UUID = evResp.Id

		return &ev, nil
	}

	return nil, err
}

