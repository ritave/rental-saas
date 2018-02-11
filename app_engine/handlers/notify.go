package handlers

import (
	"io/ioutil"
	"fmt"
	"time"
	"net/http"
	"google.golang.org/appengine"
	gae_log "google.golang.org/appengine/log"
	"calendar-synch/app_engine/logic"
)

func Notify(w http.ResponseWriter, r *http.Request) {
	//srv := GetService(r)
	ctx := appengine.NewContext(r)

	// try to discover what Google is acutally passing to us
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("Error parsing response: %s\n", err)
	}
	fmt.Printf("%s\n", string(body))
	gae_log.Infof(ctx, "Received notification %s", string(body))
	logic.RecordVisit(ctx, time.Now(), r.RemoteAddr, r.Body)
	//

	// TODO
	// Check what has changed in the calendar
	// (compare with saved version somewhere)
}

