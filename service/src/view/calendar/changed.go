package calendar

import (
	"rental-saas/service/src/model"
	"rental-saas/service/src/presenter"
	"rental-saas/service/src/application/core"
	"errors"
	"github.com/sirupsen/logrus"
)

type ChangedRequest struct{}
type ChangedResponse struct{}
type Modification struct {
	Flags []string `json:"flags"`
	model.Event
}

func Changed(a *core.Application, r interface{}) (interface{}, error) {
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
	logrus.Printf("Synchronisation had following effect: %v", effect)

	response := make([]Modification, len(diff))

	for ind, eventChanged := range diff {
		response[ind] = Modification{
			Flags: eventChanged.ToListOfWords(),
			Event: *eventChanged.Event,
		}
	}

	logrus.Debugf("Things that chagned:\n%#v", response)

	presenter.TakeActionOnDifferences(a.Utils.Pozamiatane, a.Calendar, diff)

	return ChangedResponse{}, nil
}
