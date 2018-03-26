package my_calendar

import (
	"google.golang.org/api/calendar/v3"
	"rental-saas/src/model"
	"time"
	"fmt"
	gaeLog "google.golang.org/appengine/log"
	"context"
)

// FIXME temporary
var ctx = context.Background()

var wat = func() *bool {
	b := false
	return &b
}

func AddEvent(cal *calendar.Service, ev *model.Event) (*model.Event, error){
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
		GuestsCanInviteOthers: wat(),
		GuestsCanSeeOtherGuests: wat(),
	}

	//creationDate := utils.TimeToString(time.Now())

	evResp, err := cal.Events.Insert("primary", newEvent).Do()
	if err != nil {
		gaeLog.Debugf(ctx, "Adding event failed %#v %s", ev, err.Error())
		return nil, err
	}
		//log.Printf("Link: %s", evResp.HtmlLink)
		//
		//// Creation date will be my "primary-key"
		//eventCreationTime, err := utils.VerifyStringToTime(evResp.Created)
		//if err != nil {
		//	gaeLog.Criticalf(ctx, "Google passed to us string that is not of valid format! IMPOSSIBRU!")
		//	eventCreationTime = time.Now()
		//}
		//eventOrderingKey := utils.TimeToMilliseconds(eventCreationTime)
		//ev.Timestamp = eventOrderingKey
		//ev.UUID = evResp.Id
		//
		//return &ev, nil
		return model.ConvertGoogleToMine(evResp), nil
}

