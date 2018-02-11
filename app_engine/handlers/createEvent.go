package handlers

import (
	"net/http"
	"encoding/json"
	"calendar-synch/app_engine/logic"
)

type EventRequest struct {
	Summary  string `json:"summary"`
	User     string `json:"user"`
	Start    string `json:"start"`
	End      string `json:"end"`
	Location string `json:"location"`
}

func CreateEvent(w http.ResponseWriter, r *http.Request) {
	srv := GetService(r)

	eventRequest, err := ExtractEventFromBody(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Malformed json"))
	}

	logic.AddEventToCalendar(srv, logic.Event(eventRequest))
	logic.WatchForChanges(srv)
}

func ExtractEventFromBody(r *http.Request) (EventRequest, error) {
	var target EventRequest
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&target)
	if err != nil {
		return EventRequest{}, err
	}
	return target, nil
}

