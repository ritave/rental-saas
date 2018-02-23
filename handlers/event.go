package handlers

import (
	"net/http"
	"encoding/json"
	"calendar-synch/logic"
	"google.golang.org/appengine"
	"calendar-synch/objects"
	"io/ioutil"
	"log"
)

type EventRequest struct {
	Summary      string `json:"summary"`
	User         string `json:"user"`
	Start        string `json:"start"`
	End          string `json:"end"`
	Location     string `json:"location"`
	CreationDate string `json:"-"` //not used
}

// TODO mux + contexts + jsonification of interface{}

// TODO distinction of POST, GET, OPTIONS

// TODO split this into different files for each handler

var flag = true

func EventCreate(w http.ResponseWriter, r *http.Request) {
	if flag {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8000")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	}

	if r.Method == http.MethodOptions {
		return
	}

	eventRequest, err := ExtractEventFromBody(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("\"Malformed json\""))
		return
	}

	if appengine.IsDevAppServer() {
		w.Write([]byte("\"Congratz\"")) // JSONified
		return
	}

	err = logic.EvenMoreChecksForTheEvent(objects.Event(eventRequest))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("\"" + err.Error() + "\""))
		return
	}

	srv := GetService(r)
	ctx := appengine.NewContext(r)

	// TODO move this logic level down

	// TODO rollbacks
	event, err := logic.AddEventToCalendar(srv, objects.Event(eventRequest))
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

func ExtractEventFromBody(r *http.Request) (EventRequest, error) {
	var target EventRequest
	defer r.Body.Close()

	if appengine.IsDevAppServer() {
		bytez, _ := ioutil.ReadAll(r.Body)

		log.Println("JSON received")
		log.Println(string(bytez))

		err := json.Unmarshal(bytez, &target)
		if err != nil {
			return EventRequest{}, err
		}
		return target, nil
	}

	err := json.NewDecoder(r.Body).Decode(&target)
	if err != nil {
		return EventRequest{}, err
	}
	return target, nil
}

func EventList(w http.ResponseWriter, r *http.Request) {
	if flag {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8000")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	}

	srv := GetService(r)

	events, err := srv.Events.List("primary").ShowDeleted(true).
		MaxResults(100).OrderBy("updated").Do()
	if err != nil {
		log.Println("Error fetching events")
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	result := make([]objects.Event, 0)
	for _, ev := range events.Items {
		converted, _ := logic.ConvertEventToEventLol(ev)
		result = append(result, *converted)
	}
	bytez, err := json.Marshal(&result)
	if err != nil {
		log.Println("Error marshalling response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if flag {
		log.Println("We will be sending this back:")
		log.Println(string(bytez))
	}

	w.Write(bytez)
}

/*
{"summary":"a", "user":"a", "start":"a", "end":"a", "location":"a", "creationDate":"a"}
 */

func EventChanged(w http.ResponseWriter, r *http.Request) {
	log.Println("Captain, we are being hailed.")
}
