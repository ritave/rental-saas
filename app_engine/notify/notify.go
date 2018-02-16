package main

import (
	"net/http"
	"log"
	"calendar-synch/handlers"
	"context"
)

func main() {
	http.HandleFunc("/notify/send", HandlerSend)
	log.Print("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func HandlerSend(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Ok, I got this."))
}

func init() {
	testPing()

	background := context.Background()

	srv := handlers.GetServiceWithoutRequest(background)
	events, err := srv.Events.List("primary").Do()
	if err != nil {
		log.Printf("Listing events: %s", err.Error())
	} else {
		log.Println("Upcoming events:")
		if len(events.Items) > 0 {
			for _, i := range events.Items {
				var when string
				// If the DateTime is an empty string the Event is an all-day Event.
				// So only Date is available.
				if i.Start.DateTime != "" {
					when = i.Start.DateTime
				} else {
					when = i.Start.Date
				}
				log.Printf("%s (%s)\n", i.Summary, when)
			}
		} else {
			log.Printf("No upcoming events found.\n")
		}
	}
}

func testPing() {
	resp, err := http.DefaultClient.Get("https://calendar-cron.appspot.com/event/ping")
	if err != nil {
		log.Printf("Request on init: %s", err.Error())
	} else {
		log.Printf("Response status: %d", resp.StatusCode)
	}
}