package main

import (
	"google.golang.org/api/calendar/v3"
	"net/http"
	"time"

	"log"
	"encoding/json"
)

type CreateNewEventRequest struct {
	Summary  string `json:"summary"`
	User     string `json:"user"`
	Start    string `json:"start"`
	End      string `json:"end"`
	Location string `json:"location"`
}

func CreateNewEvent(w http.ResponseWriter, r *http.Request) {
	srv := GetService(r)

	eventRequest, err := extractEventFromBody(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Malformed json"))
	}

	addSomeEventToCalendar(srv, eventRequest)
	watchForChanges(srv)
}

func extractEventFromBody(r *http.Request) (CreateNewEventRequest, error) {
	var target CreateNewEventRequest
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&target)
	if err != nil {
		return CreateNewEventRequest{}, err
	}
	return target, nil
}

func addSomeEventToCalendar(cal *calendar.Service, eventRequest CreateNewEventRequest) {
	newEvent := &calendar.Event{
		Summary:     eventRequest.Summary,
		Location:    eventRequest.Location,
		Description: "Cleaning service ordered on %s. Feel free to move this event in your calendar to change the date!",
		Start: &calendar.EventDateTime{
			DateTime: time.Now().Format(time.RFC3339),
		},
		End: &calendar.EventDateTime{
			DateTime: time.Now().Add(time.Hour).Format(time.RFC3339),
		},
		Attendees: []*calendar.EventAttendee{
			&calendar.EventAttendee{Email: eventRequest.User},
		},
		GuestsCanModify: true,
	}

	ev, err := cal.Events.Insert("primary", newEvent).Do()
	if err != nil {
		log.Println("Adding event failed")
		log.Printf("Error: %s", err.Error())
	} else {
		log.Printf("Link: %s", ev.HtmlLink)
	}
}
