package calendar

import (
	"rental-saas/src/model"
	"rental-saas/src/calendar_wrap"
	"rental-saas/src/presenter"
	"encoding/json"
	"google.golang.org/appengine/urlfetch"
	"bytes"
	"net/http"
	"log"
	"google.golang.org/appengine"
	gaeLog "google.golang.org/appengine/log"
	"rental-saas/src/presenter/my_datastore"
	"rental-saas/src/presenter/wrapper"
	"errors"
)

type ChangedResponse []Modification
type Modification struct {
	Flags []string `json:"flags"`
	model.Event
}

func Changed(a *wrapper.Application, r interface{}) (interface{}, error) {
	var err error
	request, ok := r.(ChangedResponse)
	if !ok {
		return nil, errors.New("reflection failed")
	}

	diff, err := presenter.FindChanged(ctx, cal)
	if err != nil {
		gaeLog.Debugf(ctx, "Finding changes: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// no errors returned, fingers crossed it works!
	effect := my_datastore.SynchroniseDatastore(ctx, diff)
	gaeLog.Debugf(ctx, "Synchronisation had following effect: %#v", effect)
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

	presenter.TakeActionOnDifferences(ctx, cal, diff)
}
