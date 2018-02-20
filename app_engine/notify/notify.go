package main

import (
	"net/http"
	"log"
	"calendar-synch/handlers"
	"context"
	"calendar-synch/logic"
	"google.golang.org/api/calendar/v3"
	"os"
)

const NotifyGet = "/notify/get"

const EnvAppPing = "APP_PING"
const EnvAppChanged = "APP_CHANGED"

func main() {
	http.HandleFunc(NotifyGet, HandlerGet)
	log.Print("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func HandlerGet(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Thanks Google, I got this from here."))

	notifyMainApp()
}

func init() {
	testPing()

	background := context.Background()

	cal := handlers.GetServiceWithoutRequest(background)
	registerReceiver(cal)
}

func testPing() {
	pingAddr, exists := os.LookupEnv(EnvAppPing)
	if !exists {
		pingAddr = "https://calendar-cron.appspot.com/event/ping"
		log.Printf("Resolving to prior default: %s", pingAddr)
	}

	resp, err := http.DefaultClient.Get(pingAddr)
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
	log.Println("Registering receiver")
	selfAddr, exists := os.LookupEnv(EnvAppChanged)
	if !exists {
		selfAddr = "https://calendar-cron.appspot.com/"
		log.Printf("Registering: resolving to prior default: %s", selfAddr)
	}
	err := logic.WatchForChanges(cal, selfAddr + NotifyGet)
	if err != nil {
		log.Printf("Error sending watch request: %s", err.Error())
	}
}

func notifyMainApp() {
	pingAddr, exists := os.LookupEnv(EnvAppChanged)
	if !exists {
		pingAddr = "https://calendar-cron.appspot.com/event/changed"
		log.Printf("Notify: resolving to prior default: %s", pingAddr)
	}

	resp, err := http.DefaultClient.Get(pingAddr)
	if err != nil {
		log.Printf("Notifying: %s", err.Error())
	} else {
		log.Printf("Notifying: %d", resp.StatusCode)
	}

}