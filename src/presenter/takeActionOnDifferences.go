package presenter

import (
	"rental-saas/src/model"
	"google.golang.org/api/calendar/v3"
		"rental-saas/src/application/interfaces"
	"github.com/sirupsen/logrus"
	"rental-saas/src/api_integration"
)

func TakeActionOnDifferences(pozamiatane api_integration.Provider, cal interfaces.CalendarInterface, diff []*model.EventModified) {
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
					logrus.Printf( "Hacking failed: %s", err.Error())
				}

			case model.ModifiedLocation:

			case model.ModifiedTime:
				// YOU KNOW WHAT TO DO
				// TODO lol
			}
		}
	}
}
