package event

import (
	"rental-saas/src/application/core"
	"errors"
)

type DeleteRequest struct {
	UUID string `json:"uuid"`
}

type DeleteResponse struct {
	Message string `json:"message"`
}

func Delete(a *core.Application, r interface{}) (interface{}, error) {
	var err error
	request, ok := r.(DeleteRequest)
	if !ok {
		return nil, errors.New("reflection failed")
	}

	err = a.Calendar.DeleteEvent(request.UUID)
	if err != nil {
		return nil, err
	}
	err = a.Datastore.DeleteEvent(request.UUID)
	if err != nil {
		return nil, err
	}

	return DeleteResponse{"Deleted event"}, nil
}
