package event

import (
	"rental-saas/src/model"
	"rental-saas/src/application/core"
	"errors"
)

type CreateRequest struct {
	Summary      string `json:"summary"`
	User         string `json:"user"`
	Start        string `json:"start"`
	End          string `json:"end"`
	Location     string `json:"location"`
	CreationDate string `json:"-"`
	Timestamp    int64  `json:"-"`
	UUID         string `json:"-"`
	TestFields string `json:"-"`
}

type CreateResponse struct {
	Message string `json:"message"`
}

func Create(a *core.Application, r interface{}) (interface{}, error) {
	var err error
	eventRequest, ok := r.(CreateRequest)
	if !ok {
		return nil, errors.New("reflection failed")
	}

	// FIXME
	wat := model.Event(eventRequest)
	event := &wat
	// TODO this should return 'corrected' event (even a pointer to it)
	err = model.EvenMoreChecksForTheEvent(*event)
	if err != nil {
		return nil, err
	}

	event, err = a.Calendar.AddEvent(event)
	if err != nil {
		return nil, err
	}
	err = a.Datastore.SaveEvent(event)
	if err != nil {
		return nil, err
	}

	return CreateResponse{"Created event"}, nil
}
