package main

import (
	"calendar-synch/service"
	"log"
	"fmt"
	"google.golang.org/api/calendar/v3"
	"time"
)

func main() {
	cal := service.New()

	ShowEvents(cal)

	channel := calendar.Channel{}
	cal.Events.Watch("primary", &channel)
}

func ListCalendars(cal *calendar.Service) {
	calList, err := cal.CalendarList.List().ShowDeleted(false).Do()
	if err != nil {
		log.Printf("I don't know what I'm doing actually")
		log.Printf(err.Error())
	} else {
		fmt.Println("Listing calendars:")
		for _, i := range calList.Items {
			fmt.Println(i)
		}
	}
}

func AddSomeEvent(cal *calendar.Service) {
	event := &calendar.Event{
		Summary:     "Sprzatanie",
		Location:    "800 Howard St., San Francisco, CA 94103",
		Description: "Sprzatanie jakiegos miejsca",
		Start: &calendar.EventDateTime{
			DateTime: time.Now().Format(time.RFC3339),
		},
		End: &calendar.EventDateTime{
			DateTime: time.Now().Add(time.Hour).Format(time.RFC3339),
		},
		Attendees: []*calendar.EventAttendee{
			&calendar.EventAttendee{Email: "cymerrad@gmail.com"},
		},
		GuestsCanModify: true,
	}

	ev, err := cal.Events.Insert("primary", event).Do()
	if err != nil {
		log.Println("Adding event failed")
		log.Println(err.Error())
	} else {
		log.Println("Link: ", ev.HtmlLink)
	}
}

func ShowEvents(cal *calendar.Service) {
	eventList, err := cal.Events.List("primary").Do()
	if err != nil {
		log.Printf("Listing events failed: %s", err.Error())
	} else {
		fmt.Println("Events:")
		for _, i := range eventList.Items {
			fmt.Printf("%s from %s to %s", i.Summary, i.Start.DateTime, i.End.DateTime)
	}
}

}
