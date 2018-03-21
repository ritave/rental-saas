package event

import (
	"rental-saas/src/calendar_wrap"
	"google.golang.org/appengine"
	"io/ioutil"
	"encoding/json"
	"net/http"
	"log"
	gaeLog "google.golang.org/appengine/log"
	"rental-saas/src/model/my_calendar"
	"rental-saas/src/model/my_datastore"
	"rental-saas/src/utils"
)

type DeleteRequest struct {
	UUID string `json:"uuid"`
}

func Delete(w http.ResponseWriter, r *http.Request) {
	eventRequest, err := extractDeleteRequestFromBody(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		utils.WriteAsJSON(w, "Malformed JSON")
		return
	}

	cal := calendar_wrap.NewStandard(r)
	ctx := appengine.NewContext(r)

	err = my_calendar.DeleteEvent(cal, eventRequest.UUID)
	if err != nil {
		gaeLog.Debugf(ctx, "Calendar delete %s: %s", eventRequest.UUID, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = my_datastore.DeleteEvent(ctx, eventRequest.UUID)
	if err != nil {
		gaeLog.Debugf(ctx, "Datastore delete %s: %s", eventRequest.UUID, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	utils.WriteAsJSON(w, "Delted event")
}

func extractDeleteRequestFromBody(r *http.Request) (DeleteRequest, error) {
	var target DeleteRequest
	defer r.Body.Close()

	if appengine.IsDevAppServer() {
		bytez, _ := ioutil.ReadAll(r.Body)

		log.Println("JSON received")
		log.Println(string(bytez))

		err := json.Unmarshal(bytez, &target)
		if err != nil {
			return DeleteRequest{}, err
		}
		return target, nil
	}

	err := json.NewDecoder(r.Body).Decode(&target)
	if err != nil {
		return DeleteRequest{}, err
	}
	return target, nil
}
