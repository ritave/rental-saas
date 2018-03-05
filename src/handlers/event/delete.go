package event

import (
	"calendar-synch/src/calendar_wrap"
	"google.golang.org/appengine"
	"io/ioutil"
	"encoding/json"
	"net/http"
	"log"
	"calendar-synch/src/logic/my_calendar"
	"calendar-synch/src/logic/my_datastore"
)

type DeleteRequest struct {
	UUID string `json:"uuid"`
}

func Delete(w http.ResponseWriter, r *http.Request) {
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

	eventRequest, err := extractDeleteRequestFromBody(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("\"Malformed json\""))
		return
	}

	cal := calendar_wrap.NewStandard(r)
	ctx := appengine.NewContext(r)

	err = my_calendar.DeleteEvent(cal, eventRequest.UUID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = my_datastore.DeleteEvent(ctx, eventRequest.UUID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write([]byte("\"Created event\""))
}

func extractDeleteRequestFromBody(r *http.Request) (DeleteRequest, error) {
	var target DeleteRequest
	defer r.Body.Close()

	if appengine.IsDevAppServer() {
		bytez, _ := ioutil.ReadAll(r.Body)

		log.Println("JSON received")
		log.Println(string(bytez))

		err := json.Unmarshal(bytez, &target)
		if err != nil {
			return DeleteRequest{}, err
		}
		return target, nil
	}

	err := json.NewDecoder(r.Body).Decode(&target)
	if err != nil {
		return DeleteRequest{}, err
	}
	return target, nil
}
