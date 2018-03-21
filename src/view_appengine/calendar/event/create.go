package event

import (
	"rental-saas/src/model"
	"rental-saas/src/calendar_wrap"
	"io/ioutil"
	"encoding/json"
	"google.golang.org/appengine"
	"net/http"
	"log"
	"rental-saas/src/presenter/my_calendar"
	"rental-saas/src/presenter/my_datastore"
	"rental-saas/src/utils"
	gaeLog "google.golang.org/appengine/log"
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

func Create(w http.ResponseWriter, r *http.Request) {
	eventRequest, err := extractCreateRequestFromBody(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		utils.WriteAsJSON(w, "Malformed JSON")
		return
	}

	err = model.EvenMoreChecksForTheEvent(model.Event(eventRequest))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		utils.WriteAsJSON(w, err.Error())
		return
	}

	cal := calendar_wrap.NewStandard(r)
	ctx := appengine.NewContext(r)

	event, err := my_calendar.AddEvent(ctx, cal, model.Event(eventRequest))
	if err != nil {
		gaeLog.Debugf(ctx, "Calendar create %s: %s", eventRequest.UUID, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = my_datastore.SaveEventInDatastore(ctx, event)
	if err != nil {
		gaeLog.Debugf(ctx, "Datastore create %s: %s", eventRequest.UUID, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	utils.WriteAsJSON(w, "Created event")
}

func extractCreateRequestFromBody(r *http.Request) (CreateRequest, error) {
	var target CreateRequest
	defer r.Body.Close()

	if appengine.IsDevAppServer() {
		bytez, _ := ioutil.ReadAll(r.Body)

		log.Println("JSON received")
		log.Println(string(bytez))

		err := json.Unmarshal(bytez, &target)
		if err != nil {
			return CreateRequest{}, err
		}
		return target, nil
	}

	err := json.NewDecoder(r.Body).Decode(&target)
	if err != nil {
		return CreateRequest{}, err
	}
	return target, nil
}
