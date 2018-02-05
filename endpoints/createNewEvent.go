package endpoints

import (
	"google.golang.org/api/calendar/v3"
	"net/http"
	"time"
	"fmt"

	"log"
)

type CreateNewEventRequest struct {
	Summary  string `json:"summary"`
	User     string `json:"user"`
	Start    string `json:"start"`
	End      string `json:"end"`
	Location string `json:"location"`
}

func CreateNewEvent(s *calendar.Service, w http.ResponseWriter, r *http.Request) {

	msg := fmt.Sprintf("Added from the cloud by %S", s.UserAgent)
	AddSomeEventToCalendar(s, msg)

}

func AddSomeEventToCalendar(cal *calendar.Service, summary string) {
	event := &calendar.Event{
		Summary:     summary,
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
		log.Printf("Error: %s", err.Error())
	} else {
		log.Printf("Link: %s", ev.HtmlLink)
	}
}
