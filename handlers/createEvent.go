package handlers

import (
	"net/http"
	"encoding/json"
	"calendar-synch/logic"
	"google.golang.org/appengine"
	"calendar-synch/objects"
	"calendar-synch/helpers"
)

type EventRequest struct {
	Summary      string `json:"summary"`
	User         string `json:"user"`
	Start        string `json:"start"`
	End          string `json:"end"`
	Location     string `json:"location"`
	CreationDate string `json:"-"` //not used
}

func CreateEvent(w http.ResponseWriter, r *http.Request) {
	defer helpers.RecoverPanic()

	srv := GetService(r)
	ctx := appengine.NewContext(r)

	eventRequest, err := ExtractEventFromBody(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Malformed json"))
	}

	// TODO move this logic level down
	event, err := logic.AddEventToCalendar(srv, objects.Event(eventRequest))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	err = logic.SaveEventInDatastore(ctx, event)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	// TODO this should be called at best only once... Not at every CreateEvent call.
	// TODO also there are some refreshing tokens flying around soo... yeeeah...
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
