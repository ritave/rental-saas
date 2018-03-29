package calendar

import (
	"rental-saas/src/model"
	"rental-saas/src/presenter"
	"encoding/json"
	"bytes"
	"net/http"
	"log"
	"rental-saas/src/presenter/wrapper"
	"errors"
)

type ChangedRequest struct{}
type ChangedResponse struct{}
type Modification struct {
	Flags []string `json:"flags"`
	model.Event
}

func Changed(a *wrapper.Application, r interface{}) (interface{}, error) {
	var err error
	_, ok := r.(ChangedRequest)
	if !ok {
		return nil, errors.New("reflection failed")
	}

	diff, err := presenter.FindChanged(a.Datastore, a.Calendar)
	if err != nil {
		return nil, err
	}

	// no errors returned, fingers crossed it works!
	effect := a.Datastore.SynchroniseDatastore(diff)
	log.Printf("Synchronisation had following effect: %v", effect)

	response := make([]Modification, len(diff))

	for ind, eventChanged := range diff {
		response[ind] = Modification{
			Flags: eventChanged.ToListOfWords(),
			Event: *eventChanged.Event,
		}
	}

	bytez, err := json.Marshal(&response)
	if err != nil {
		return nil, err
	}

	whereTo := "https://calendarcron.appspot.com/dummy/send"

	client := http.DefaultClient
	resp, err := client.Post(whereTo, "application/json", bytes.NewReader(bytez))
	if err != nil {
		log.Printf("Error sending changes to %s: %s", whereTo, err.Error())
	} else {
		log.Println("Success sending that son of a bitch")
		log.Println(*resp)
	}

	presenter.TakeActionOnDifferences(a.Calendar, diff)

	return ChangedResponse{}, nil
}
