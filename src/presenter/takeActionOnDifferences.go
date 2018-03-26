package presenter

import (
	"rental-saas/src/model"
	"google.golang.org/api/calendar/v3"
	"log"
	"rental-saas/src/presenter/interfaces"
)

func TakeActionOnDifferences(cal interfaces.CalendarInterface, diff []*model.EventModified) {
	for _, event := range diff {
		for k := range event.Modifications {
			switch k {
			case model.Deleted:
				// TODO send again, with instructions on how to delete this
				err := cal.UpdateEvent(event.Event.UUID, &calendar.Event{
					Attendees: []*calendar.EventAttendee{{Email: event.Event.User}},
					Description: "In order to PROPERLY delete the event visit this link [TO BE ADDED]",
				})

				if err != nil {
					log.Printf( "Hacking failed: %s", err.Error())
				}

			case model.ModifiedLocation:
				// YOU KNOW WHAT TO DO
				// TODO lol
			case model.ModifiedTime:
				// YOU KNOW WHAT TO DO
				// TODO lol
			}
		}
	}
}
