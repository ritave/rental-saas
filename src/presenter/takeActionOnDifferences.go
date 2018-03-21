package presenter

import (
	"rental-saas/src/model"
	"google.golang.org/api/calendar/v3"
	"context"
	gaeLog "google.golang.org/appengine/log"
)

func TakeActionOnDifferences(ctx context.Context, cal *calendar.Service, diff []*model.EventModified) {
	for _, event := range diff {
		for k := range event.Modifications {
			switch k {
			case model.Deleted:
				// send again, with instructions on how to delete this
				resp, err := cal.Events.Update("primary", event.Event.UUID, &calendar.Event{
					Attendees: []*calendar.EventAttendee{{Email: event.Event.User}},
					Description: "In order to PROPERLY delete the event visit this link [TO BE ADDED]",
				}).Do()
				if err != nil {
					gaeLog.Debugf(ctx, "Hacking failed: %s", err.Error())
				} else {
					gaeLog.Debugf(ctx, "Hacking succeeded: %#v", resp)
				}
			case model.ModifiedLocation:
				// YOU KNOW WHAT TO DO
			case model.ModifiedTime:
				// YOU KNOW WHAT TO DO
			}
		}
	}
}
