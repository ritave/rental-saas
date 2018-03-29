package calendar

import (
	"rental-saas/src/model"
	"rental-saas/src/presenter/wrapper"
)

type ViewRequest struct{}
type ViewResponse []*model.Event

func View(a *wrapper.Application, r interface{}) (interface{}, error) {
	//events, err := srv.Events.List("primary").ShowDeleted(false).OrderBy("updated").Do()
	return a.Calendar.QueryEvents()
}

