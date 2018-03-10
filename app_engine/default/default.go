package main

import (
	"net/http"
	"google.golang.org/appengine"
	"log"
	"calendar-synch/src/handlers/calendar"
	"calendar-synch/src/handlers/calendar/event"
)

func main() {
	bindEndpoints()
	appengine.Main()
}

func bindEndpoints() {
	http.HandleFunc("/calendar/event/create", event.Create)
	http.HandleFunc("/calendar/event/delete", event.Delete)

	http.HandleFunc("/calendar/changed", calendar.Changed)
	http.HandleFunc("/calendar/list", calendar.View)

	log.Println("Bound endpoints.")
}

