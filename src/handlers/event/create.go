package event

import (
	"calendar-synch/src/logic"
	"calendar-synch/src/objects"
	"calendar-synch/src/calendar_wrap"
	"io/ioutil"
	"encoding/json"
	"google.golang.org/appengine"
	"net/http"
	"log"
)

type CreateRequest struct {
	Summary      string `json:"summary"`
	User         string `json:"user"`
	Start        string `json:"start"`
	End          string `json:"end"`
	Location     string `json:"location"`
	Timestamp    int64  `json:"-"` //not used
	CreationDate string `json:"-"`
	UUID         string `json:"-"`
}

// TODO mux + contexts + jsonification of interface{}

// TODO distinction of POST, GET, OPTIONS

// TODO split this into different files for each handler

// TODO -> env var
var allowAccessFromLocalhost = true

const (
	CORSlocalhost = "http://localhost:8000"
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

	err = logic.EvenMoreChecksForTheEvent(objects.Event(eventRequest))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("\"" + err.Error() + "\""))
		return
	}

	if dev {
		w.Write([]byte("\"Congratz\"")) // JSONified
		return
	}

	cal := calendar_wrap.NewStandard(r)
	ctx := appengine.NewContext(r)

	// TODO move this logic level down

	// TODO rollbacks
	event, err := logic.AddEventToCalendar(cal, objects.Event(eventRequest))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = logic.SaveEventInDatastore(ctx, event)
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
