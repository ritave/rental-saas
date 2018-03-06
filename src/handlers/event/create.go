package event

import (
	"calendar-synch/src/objects"
	"calendar-synch/src/calendar_wrap"
	"io/ioutil"
	"encoding/json"
	"google.golang.org/appengine"
	"net/http"
	"log"
	"calendar-synch/src/logic/my_calendar"
	"calendar-synch/src/logic/my_datastore"
)

type CreateRequest struct {
	Summary      string `json:"summary"`
	User         string `json:"user"`
	Start        string `json:"start"`
	End          string `json:"end"`
	Location     string `json:"location"`
	CreationDate string `json:"-"`
	Timestamp    int64  `json:"-"` //not used
	UUID         string `json:"-"`
}

// TODO mux + contexts + jsonification of interface{}

// TODO distinction of POST, GET, OPTIONS

// TODO split this into different files for each handler

// TODO -> env var
var allowAccessFromLocalhost = true

const (
	CORSlocalhost = "http://localhost:8080"
	CORSapp       = "https://calendarcron.appspot.com"
)

var dev = appengine.IsDevAppServer()

func Create(w http.ResponseWriter, r *http.Request) {
	if allowAccessFromLocalhost {
		if appengine.IsDevAppServer() {
			w.Header().Set("Access-Control-Allow-Origin", CORSlocalhost)
		} else {
			w.Header().Set("Access-Control-Allow-Origin", CORSapp)
		}
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	}

	if r.Method == http.MethodOptions {
		return
	}

	eventRequest, err := extractCreateRequestFromBody(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("\"Malformed json\""))
		return
	}

	err = objects.EvenMoreChecksForTheEvent(objects.Event(eventRequest))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("\"" + err.Error() + "\""))
		return
	}

	cal := calendar_wrap.NewStandard(r)
	ctx := appengine.NewContext(r)

	// TODO move this logic level down

	// TODO rollbacks
	event, err := my_calendar.AddEvent(cal, objects.Event(eventRequest))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = my_datastore.SaveEventInDatastore(ctx, event)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write([]byte("\"Created event\""))
}

func extractCreateRequestFromBody(r *http.Request) (CreateRequest, error) {
	var target CreateRequest
	defer r.Body.Close()

	if appengine.IsDevAppServer() {
		bytez, _ := ioutil.ReadAll(r.Body)

		log.Println("JSON received")
		log.Println(string(bytez))

		err := json.Unmarshal(bytez, &target)
		if err != nil {
			return CreateRequest{}, err
		}
		return target, nil
	}

	err := json.NewDecoder(r.Body).Decode(&target)
	if err != nil {
		return CreateRequest{}, err
	}
	return target, nil
}
