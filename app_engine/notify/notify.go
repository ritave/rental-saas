package main

import (
	"net/http"
	"log"
	"calendar-synch/handlers"
	"context"
	"calendar-synch/logic"
	"google.golang.org/api/calendar/v3"
	"os"
	"strconv"
	"time"
	"calendar-synch/helpers"
)

const NotifyGet = "/notify/get"

const (
	EnvAppPing = "CALENDAR_APP_PING"
	EnvAppChanged = "CALENDAR_APP_CHANGED"
	EnvApp = "CALENDAR_APP"

	NotifyExpireAfter = "NOTIFY_EXPIRE_AFTER"
)


func main() {
	http.HandleFunc(NotifyGet, HandlerGet)
	log.Print("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func HandlerGet(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Thanks Google, I got this from here."))

	err := notifyMainApp()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func init() {
	testPing()

	background := context.Background()
	cal := handlers.GetServiceWithoutRequest(background)
	registerReceiver(cal)
}

func testPing() {
	pingAddr := getStringFromEnv(EnvAppPing, "https://calendar-cron.appspot.com/event/ping")

	resp, err := http.DefaultClient.Get(pingAddr)
	if err != nil {
		log.Printf("Request on init: %s", err.Error())
	} else {
		log.Printf("Response status: %d", resp.StatusCode)
	}
}

func registerReceiver(cal *calendar.Service) {
	log.Println("Registering receiver")
	selfAddr := getStringFromEnv(EnvApp, "https://calendar-cron.appspot.com/")

	// TODO refresh after every some constant time interval?
	// TODO also there are some refreshing tokens flying around soo... yeeeah...
	
	expireAfter, err := strconv.Atoi(getStringFromEnv(NotifyExpireAfter, "3600"))
	if err != nil {
		log.Fatalf("ATOI: %s", err.Error())
	}

	err = logic.WatchForChanges(cal, selfAddr + NotifyGet, time.Duration(expireAfter)*time.Second)
	if err != nil {
		log.Printf("Error sending watch request: %s", err.Error())

		go func() {
			log.Printf("Retrying in one minute")
			timer := time.NewTimer(time.Duration(time.Minute))
			refreshTime := <- timer.C
			log.Printf("Refreshing watch channel on %s", refreshTime.Format(helpers.DefaultTimeType))
			registerReceiver(cal)
		}()
		return
	}

	// if everything went smoothly, carry on with usual refresh-after-an-hour-or-so
	go func() {
		log.Printf("Scheduled for retrying in %d second", expireAfter)
		timer := time.NewTimer(time.Duration(expireAfter)*time.Second)
		refreshTime := <- timer.C
		log.Printf("Refreshing watch channel on %s", refreshTime.Format(helpers.DefaultTimeType))
		registerReceiver(cal)
	}()
}

func notifyMainApp() (error) {
	notifyAddr := getStringFromEnv(EnvAppChanged, "https://calendar-cron.appspot.com/event/changed")

	resp, err := http.DefaultClient.Get(notifyAddr)
	if err != nil {
		log.Printf("Notifying: %s", err.Error())
	} else {
		log.Printf("Notifying: %d", resp.StatusCode)
	}

	return err
}

func getStringFromEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Printf("Env var %s not found; fallback to: %s", key, fallback)
		return fallback
	}
	return value
}