package event

import (
	"rental-saas/src/model"
	"rental-saas/src/calendar_wrap"
	"rental-saas/src/presenter/my_calendar"
	"rental-saas/src/presenter/my_datastore"
	"rental-saas/src/presenter/wrapper"
	"errors"
	"google.golang.org/appengine"
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

func Create(a *wrapper.Application, r interface{}) (interface{}, error) {
	var err error
	eventRequest, ok := r.(CreateRequest)
	if !ok {
		return nil, errors.New("reflection failed")
	}

	err = model.EvenMoreChecksForTheEvent(model.Event(eventRequest))
	if err != nil {
		return nil, err
	}

	cal := calendar_wrap.NewStandard(r)
	ctx := appengine.NewContext(r)

	event, err := my_calendar.AddEvent(ctx, cal, model.Event(eventRequest))
	if err != nil {
		return nil, err
	}
	err = my_datastore.SaveEventInDatastore(ctx, event)
	if err != nil {
		return nil, err
	}

	return CreateResponse{"Created event"}, nil
}
