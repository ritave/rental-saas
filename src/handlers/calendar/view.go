package calendar

import (
	"calendar-synch/src/calendar_wrap"
	"calendar-synch/src/objects"
	"calendar-synch/src/logic"
	"encoding/json"
	"net/http"
	"google.golang.org/appengine"
	"log"
)

var allowAccessFromLocalhost = true

const (
	CORSlocalhost = "http://localhost:5000"
	CORSapp       = "https://calendarcron.appspot.com"
)
func View(w http.ResponseWriter, r *http.Request) {
	if allowAccessFromLocalhost {
		if appengine.IsDevAppServer() {
			w.Header().Set("Access-Control-Allow-Origin", CORSlocalhost)
		} else {
			w.Header().Set("Access-Control-Allow-Origin", CORSapp)
		}
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	}

	srv := calendar_wrap.NewStandard(r)

	events, err := srv.Events.List("primary").ShowDeleted(false).OrderBy("updated").Do()
	if err != nil {
		log.Println("Error fetching events")
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	result := make([]objects.Event, 0)
	for _, ev := range events.Items {
		converted, _ := logic.ConvertGoogleEventToMyEvent(ev)
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

