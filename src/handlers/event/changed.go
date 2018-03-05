package event

import (
	"calendar-synch/src/objects"
	"calendar-synch/src/calendar_wrap"
	"calendar-synch/src/logic"
	"encoding/json"
	"google.golang.org/appengine/urlfetch"
	"bytes"
	"net/http"
	"log"
	"google.golang.org/appengine"
	gaeLog "google.golang.org/appengine/log"
	"calendar-synch/src/logic/my_datastore"
)

type ChangedResponse []Modification
type Modification struct {
	Flags []string `json:"flags"`
	objects.Event
}

func Changed(w http.ResponseWriter, r *http.Request) {
	log.Println("Captain, we are being hailed.")

	srv := calendar_wrap.NewStandard(r)
	ctx := appengine.NewContext(r)

	diff, err := logic.FindChanged(ctx, srv)
	if err != nil {
		gaeLog.Debugf(ctx, "Finding changes: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// no errors returned, fingers crossed it works!
	effect := my_datastore.SynchroniseDatastore(ctx, diff)
	gaeLog.Debugf(ctx, "Synchronisation had following effect: %v", effect)
	//log.Printf("Synchronisation had following effect: %v", effect)

	response := make([]Modification, len(diff))

	for ind, eventChanged := range diff {
		response[ind] = Modification{
			Flags: eventChanged.ToListOfWords(),
			Event: *eventChanged.Event,
		}
	}

	bytez, err := json.Marshal(&response)
	if err != nil {
		log.Println("Error parsing response in Changed:", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	whereTo := "https://calendarcron.appspot.com/dummy/send"
	if appengine.IsDevAppServer() {
		log.Println("Replying with this response BACK to the source")
		w.Write(bytez)
		return
	}

	// X-Appengine-Inbound-Appid ?

	client := urlfetch.Client(ctx)
	resp, err := client.Post(whereTo, "application/json", bytes.NewReader(bytez))
	if err != nil {
		log.Printf("Error sending changes to %s: %s", whereTo, err.Error())
		gaeLog.Debugf(ctx, "Error sending changes to %s: %s", whereTo, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		log.Println("Unlikely success sending that son of a bitch")
		log.Println(*resp)
	}
}
