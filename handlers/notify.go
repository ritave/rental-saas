package handlers

import (
	"time"
	"net/http"
	"google.golang.org/appengine"
	"calendar-synch/logic"
	gae_log "google.golang.org/appengine/log"
	"calendar-synch/objects"
)

func Notify(w http.ResponseWriter, r *http.Request) {
	srv := GetService(r)
	ctx := appengine.NewContext(r)

	// try to discover what Google is acutally passing to us
	objects.RecordVisit(ctx, time.Now(), r.RemoteAddr, r.Body)
	//

	// Check what has changed in the calendar
	// (compare with saved version somewhere)
	difference, err := logic.FindChanged(ctx, srv)
	if err != nil {
		gae_log.Debugf(ctx, "Finding differencce: %s", err.Error())
	} else {
		gae_log.Infof(ctx, "Difference %v", difference)
	}
}

