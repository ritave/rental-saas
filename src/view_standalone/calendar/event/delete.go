package event

import (
	"rental-saas/src/presenter/wrapper"
	"errors"
)

type DeleteRequest struct {
	UUID string `json:"uuid"`
}

type DeleteResponse struct {
	Message string `json:"message"`
}

func Delete(a *wrapper.Application, r interface{}) (interface{}, error) {
	var err error
	request, ok := r.(DeleteRequest)
	if !ok {
		return nil, errors.New("reflection failed")
	}

	err = a.Calendar.DeleteEvent(request.UUID)
	if err != nil {
		return nil, err
	}
	err = a.DB.DeleteEvent(request.UUID)
	if err != nil {
		return nil, err
	}

	return DeleteResponse{"Deleted event"}, nil
}
