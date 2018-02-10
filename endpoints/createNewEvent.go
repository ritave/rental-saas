package endpoints

import (
	"google.golang.org/api/calendar/v3"
	"net/http"
	"time"
	"fmt"

	"log"
	"golang.org/x/oauth2/google"
	"google.golang.org/appengine"
	"io/ioutil"
)

type CreateNewEventRequest struct {
	Summary  string `json:"summary"`
	User     string `json:"user"`
	Start    string `json:"start"`
	End      string `json:"end"`
	Location string `json:"location"`
}

func GetServiceAccount(r *http.Request) *calendar.Service {
	if !appengine.IsDevAppServer() {
		client, err := google.DefaultClient(appengine.NewContext(r), calendar.CalendarScope)
		if err != nil {
			log.Fatalf("Default client failed: %s", err.Error())
		}
		srv, err := calendar.New(client)
		if err != nil {
			log.Fatalf("Creating calendar service on the spot failed: %s", err.Error())
		}
		return srv
	} else {
		b, err := ioutil.ReadFile("secrets/service_client_default.json")
		if err != nil {
			log.Fatalf("Unable to read client secret file: %v", err)
		}

		config, err := google.JWTConfigFromJSON(b, calendar.CalendarScope)
		if err != nil {
			log.Fatalf("Unable to parse service client secret file to config: %v", err)
		}
		client := config.Client(appengine.NewContext(r))

		srv, err := calendar.New(client)
		if err != nil {
			log.Fatalf("Unable to retrieve calendar Client %v", err)
		}
		return srv
	}
}

func CreateNewEvent(w http.ResponseWriter, r *http.Request) {
	srv := GetServiceAccount(r)

	msg := fmt.Sprintf("Added from the cloud by %s", srv.UserAgent)
	AddSomeEventToCalendar(srv, msg)
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
