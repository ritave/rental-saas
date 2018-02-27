package calendar

import (
	"net/http"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"html/template"
	"calendar-synch/src/objects"
	"calendar-synch/src/logic"
	"calendar-synch/src/calendar_wrap"
)

var tpl = template.Must(template.ParseGlob("templates/*.html"))

func View(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "GET requests only", http.StatusMethodNotAllowed)
		return
	}

	ctx := appengine.NewContext(r)
	cal := calendar_wrap.NewStandard(r)


	var currentEvents []*objects.Event
	events, err := cal.Events.List("primary").ShowDeleted(false).
		SingleEvents(true).OrderBy("startTime").Do()
	if err != nil {
		log.Debugf(ctx, "Unable to retrieve next ten of the user's events. %v", err)
	}
	currentEvents = logic.EventsMap(events.Items, logic.ConvertGoogleEventToMyEvent)

	log.Infof(ctx, "Rendering %d events", len(currentEvents))

	if err := tpl.ExecuteTemplate(w, "calendarView.html", currentEvents); err != nil {
		log.Errorf(ctx, "%v", err)
	}
}
