package handlers

import (
	"net/http"
	"encoding/json"
	"calendar-synch/logic"
	"google.golang.org/appengine"
	"calendar-synch/objects"
	"io/ioutil"
	"log"
	"bytes"
)

type EventCreateRequest struct {
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

// TODO -> env var
var allowAccessFromLocalhost = true
const (
	CORSlocalhost = "http://localhost:8000"
	CORSapp = "https://calendarcron.appspot.com"
)
var dev = appengine.IsDevAppServer()

func EventCreate(w http.ResponseWriter, r *http.Request) {
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

	eventRequest, err := ExtractEventFromBody(r)
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

	cal := GetCalendar(r)
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

func ExtractEventFromBody(r *http.Request) (EventCreateRequest, error) {
	var target EventCreateRequest
	defer r.Body.Close()

	if appengine.IsDevAppServer() {
		bytez, _ := ioutil.ReadAll(r.Body)

		log.Println("JSON received")
		log.Println(string(bytez))

		err := json.Unmarshal(bytez, &target)
		if err != nil {
			return EventCreateRequest{}, err
		}
		return target, nil
	}

	err := json.NewDecoder(r.Body).Decode(&target)
	if err != nil {
		return EventCreateRequest{}, err
	}
	return target, nil
}

func EventList(w http.ResponseWriter, r *http.Request) {
	if allowAccessFromLocalhost {
		if appengine.IsDevAppServer() {
			w.Header().Set("Access-Control-Allow-Origin", CORSlocalhost)
		} else {
			w.Header().Set("Access-Control-Allow-Origin", CORSapp)
		}
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	}

	srv := GetCalendar(r)

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

	if appengine.IsDevAppServer() {
		log.Println("We will be sending this back:")
		log.Println(string(bytez))
	}

	w.Write(bytez)
}

/*
{"summary":"a", "user":"a", "start":"a", "end":"a", "location":"a", "creationDate":"a"}
 */

type EventChangedResponse []EventModification
type EventModification struct {
	Modification []string `json:"modifications"`
	objects.Event
}

func EventChanged(w http.ResponseWriter, r *http.Request) {
	log.Println("Captain, we are being hailed.")

	srv := GetCalendar(r)
	ctx := appengine.NewContext(r)

	diff, err := logic.FindChanged(ctx, srv)
	if err != nil {
		return
	}

	response := make([]EventModification, len(diff))

	for ind, eventChanged := range diff {
		response[ind] = EventModification{
			Modification: eventChanged.ToListOfWords(),
			Event: *eventChanged.Event,
		}
	}

	bytez, err := json.Marshal(&response)
	if err != nil {
		log.Println("Error parsing response in EventChanged:", err.Error())
		return
	}


	whereTo := "https://calendarcron.appspot.com/dummy/send"
	if appengine.IsDevAppServer() {
		whereTo = "http://localhost:8081" // TODO will it be really that?
	}

	resp, err := http.DefaultClient.Post(whereTo, "application/json", bytes.NewReader(bytez))
	if err != nil {
		log.Printf("Error sending changes to %s: %s", whereTo, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		log.Println("Unlikely success sending that son of a bitch")
		log.Println(*resp)
	}
}
