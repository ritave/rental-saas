package main

import (
	"net/http"
	"log"
	"context"
	"rental-saas/src/presenter"
	"google.golang.org/api/calendar/v3"
	"os"
	"strconv"
	"time"
	"rental-saas/src/utils"
	"encoding/json"
	"rental-saas/src/calendar_wrap"
	"rental-saas/src/view_appengine/notify"
)

const NotifyGet = "/notify/get"
const NotifyPing = "/notify/ping"
const NotifyChannelDelete = "/notify/channel/delete"

const (
	EnvAppChanged = "CALENDAR_APP_CHANGED"
	EnvApp = "CALENDAR_APP"

	NotifyExpireAfter = "NOTIFY_EXPIRE_AFTER"
)

var lastReceipt presenter.ImportantChannelFields
var ticker *utils.Ticker

func main() {
	http.HandleFunc(NotifyGet, HandlerGet)
	http.HandleFunc(NotifyPing, HandlerPing)
	http.HandleFunc(NotifyChannelDelete, notify.DeleteChannel)

	log.Print("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func HandlerGet(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Thanks Google, I got this from here."))

	ticker.Restart()
}

var emptyICF = presenter.ImportantChannelFields{}
func HandlerPing(w http.ResponseWriter, r *http.Request) {
	if lastReceipt == emptyICF {
		w.Write([]byte("{\"Error\":\"there is no receiver registered\"}"))
	} else {
		bytez, _ := json.Marshal(&struct{
			Error string
			presenter.ImportantChannelFields
		}{
			Error: "congratz",
			ImportantChannelFields: lastReceipt,
		})
		w.Write(bytez)
	}
}

func init() {
	background := context.Background()
	cal := calendar_wrap.NewFlex(background)

	if cal == nil {
		log.Fatalf("Calendar As A Service was a nil")
	}

	ticker = utils.New(3*time.Second, func(){
		err := notifyMainApp()
		if err != nil {
			log.Printf("Notifying error: %s", err.Error())
		} else {
			log.Printf("Successful notifying")
		}
	})

	registerReceiver(cal)
	err := notifyMainApp()
	if err != nil {
		log.Printf("Notifying at init failed %s", err.Error())
	}
}

func registerReceiver(cal *calendar.Service) {
	log.Println("Registering receiver")
	selfAddr := getStringFromEnv(EnvApp, "https://calendarcron.appspot.com/")

	// TODO refresh after every some constant time interval?
	// TODO also there are some refreshing tokens flying around soo... yeeeah...
	
	expireAfter, err := strconv.Atoi(getStringFromEnv(NotifyExpireAfter, "3600"))
	if err != nil {
		log.Fatalf("ATOI: %s", err.Error())
	}

	err, channelReceipt := presenter.WatchForChanges(cal, selfAddr + NotifyGet, time.Duration(expireAfter)*time.Second)
	if err != nil {
		log.Printf("Error sending watch request: %s", err.Error())

		go func() {
			log.Printf("Retrying in one minute")
			timer := time.NewTimer(time.Duration(time.Minute))
			refreshTime := <- timer.C
			log.Printf("Refreshing watch channel on %s", refreshTime.Format(utils.DefaultTimeType))
			registerReceiver(cal)
		}()
		return
	}

	// global thingy for pingy
	lastReceipt = channelReceipt
	// wat

	// if everything went smoothly, carry on with usual refresh-after-an-hour-or-so
	go func() {
		log.Printf("Scheduled for retrying in %d second", expireAfter)
		timer := time.NewTimer(time.Duration(expireAfter)*time.Second)
		refreshTime := <- timer.C
		log.Printf("Refreshing watch channel on %s", refreshTime.Format(utils.DefaultTimeType))
		registerReceiver(cal)
	}()
}

func notifyMainApp() (error) {
	notifyAddr := getStringFromEnv(EnvAppChanged, "https://calendarcron.appspot.com/event/changed")

	_, err := http.DefaultClient.Get(notifyAddr)
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