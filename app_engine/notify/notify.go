package main

import (
	"net/http"
	"log"
	"calendar-synch/handlers"
	"context"
	"calendar-synch/logic"
	"google.golang.org/api/calendar/v3"
)

const NotifyGet = "/notify/get"

func main() {
	http.HandleFunc(NotifyGet, HandlerGet)
	log.Print("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func HandlerGet(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Ok, I got this."))
}

func init() {
	testPing()

	background := context.Background()

	cal := handlers.GetServiceWithoutRequest(background)
	registerReceiver(cal)
}

func testPing() {
	resp, err := http.DefaultClient.Get("https://calendar-cron.appspot.com/event/ping")
	if err != nil {
		log.Printf("Request on init: %s", err.Error())
	} else {
		log.Printf("Response status: %d", resp.StatusCode)
	}
}

func registerReceiver(cal *calendar.Service) {
	// TODO this should be called at best only once...
	// TODO also there are some refreshing tokens flying around soo... yeeeah...

	// TODO error handling

	// TODO refresh after every some constant time interval?
	logic.WatchForChanges(cal, "https://calendar-cron.appspot.com" + NotifyGet)
}