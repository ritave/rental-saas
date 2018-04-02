package presenter

import (
	"rental-saas/src/model"
	"google.golang.org/api/calendar/v3"
		"rental-saas/src/application/interfaces"
	"github.com/sirupsen/logrus"
	"rental-saas/src/api_integration"
	"rental-saas/src/utils"
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
				userID := event.Event.UserID
				orderID := event.Event.OrderID

				cleaningEdit := api_integration.EditRequest{
					UserID: userID,
					OrderID: orderID,
					Address: api_integration.Address{
						Street: event.Event.Location, // FIXME well fuck...
					},
				}

				req, err := pozamiatane.NewRequest(api_integration.EditAction, cleaningEdit)
				if err != nil {
					logrus.Printf("Modified location: %s", err.Error())
				}

				pozamiatane.SendRequestJustLog(req)

			case model.ModifiedTime:
				userID := event.Event.UserID
				orderID := event.Event.OrderID

				start := utils.StringToTime(event.Event.Start)
				end := utils.StringToTime(event.Event.End)
				cleaningDuration := end.Sub(start).Hours() // fingers-fucking-crossed

				cleaningEdit := api_integration.EditRequest{
					UserID: userID,
					OrderID: orderID,
					CleaningDate: utils.POZAMIATANE_DatetimeToDateString(start), // this SHOULD work...
					CleaningTime: utils.POZAMIATANE_DatetimeToTimeString(start),
					Length: cleaningDuration,
				}

				req, err := pozamiatane.NewRequest(api_integration.EditAction, cleaningEdit)
				if err != nil {
					logrus.Printf("Modified time: %s", err.Error())
				}

				pozamiatane.SendRequestJustLog(req)
			}
		}
	}
}
