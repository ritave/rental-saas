package calendar

import (
	"calendar-synch/src/calendar_wrap"
	"calendar-synch/src/objects"
	"encoding/json"
	"net/http"
	"google.golang.org/appengine"
	"log"
	gaeLog "google.golang.org/appengine/log"
)

func View(w http.ResponseWriter, r *http.Request) {
	srv := calendar_wrap.NewStandard(r)
	ctx := appengine.NewContext(r)

	events, err := srv.Events.List("primary").ShowDeleted(false).OrderBy("updated").Do()
	if err != nil {
		gaeLog.Debugf(ctx, "Listing events %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	result := make([]objects.Event, 0)
	for _, ev := range events.Items {
		converted, _ := objects.ConvertGoogleToMine(ev)
		result = append(result, *converted)
	}
	bytez, err := json.Marshal(&result)
	if err != nil {
		gaeLog.Debugf(ctx, "Marshalling events %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if appengine.IsDevAppServer() {
		log.Println("We will be sending this back:")
		log.Println(string(bytez))
	}

	w.Write(bytez)
}

