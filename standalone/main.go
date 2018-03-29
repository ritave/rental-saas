package main

import (
	"net/http"
	"log"
	"time"
	"rental-saas/src/utils"
	"rental-saas/src/view/notify"
	"github.com/rs/cors"
	"rental-saas/src/view/calendar/event"
	"rental-saas/src/view/calendar"
	"fmt"
	. "rental-saas/src/presenter/wrapper"
	"rental-saas/src/utils/config"
	"rental-saas/src/presenter/interfaces"
)

const NotifyGet = "/notify/get"

const (
	EnvAppChanged = "CALENDAR_APP_CHANGED"
	EnvApp = "CALENDAR_APP"

	NotifyExpireAfter = "NOTIFY_EXPIRE_AFTER"
)

var ticker *utils.Ticker

func main() {
	conf := config.GetConfig()
	app := New(conf)
	mux := http.NewServeMux()

	// events related
	mux.Handle("/calendar/event/create", &AppHandler{app, event.CreateRequest{}, event.Create})
	mux.Handle("/calendar/event/delete", &AppHandler{app, event.DeleteRequest{}, event.Delete})

	// calendar related
	mux.Handle("/calendar/changed", &AppHandler{app, calendar.ChangedRequest{}, calendar.Changed})
	mux.Handle("/calendar/view", &AppHandler{app, calendar.ViewRequest{}, calendar.View})
	
	// notify related
	mux.HandleFunc("/notify/get", HandlerGet)
	mux.HandleFunc("/notify/channel/delete", notify.DeleteChannel)

	// keep alive & admin retarted
	mux.HandleFunc("/ping", HandlerPing)

	// cors
	c := cors.New(cors.Options{
		AllowedOrigins: conf.CORS,
	})
	handler := c.Handler(mux)
	http.Handle("/", handler)

	if conf.Db.Restart {
		log.Println("Dropping and creating tables in database")
		app.Datastore.Restart()
	}

	notifySetup(app)

	// TODO this would make a good test (sans getting events from calendar)
	//events, err := app.Calendar.QueryEvents()
	//if err != nil {
	//	log.Fatalf("Querying from calendar: %s", err.Error())
	//}
	//log.Printf("Brought %d events", len(events))
	//for _, ev := range events {
	//	err := app.Datastore.SaveEvent(ev)
	//	if err != nil {
	//		log.Printf("Putting event: %s", err.Error())
	//		log.Printf("Culprit: %#v", *ev)
	//	}
	//}
	//
	//events, err = app.Datastore.QueryEvents()
	//if err != nil {
	//	log.Fatalf("Bringing events from the dead: %s", err.Error())
	//}
	//log.Printf("Brought %d events", len(events))

	log.Print("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}


// FIXME
// some ugly fuckery going on down here

func HandlerGet(w http.ResponseWriter, r *http.Request) {
	ticker.Restart()
}

func HandlerPing(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Pong")
}

func notifySetup(application *Application) {
	notifyAddr := application.Config.Receiver.Channel
	expireAfter := application.Config.Receiver.Expiration

	ticker = utils.New(time.Duration(application.Config.Receiver.Expiration)*time.Second, func(){
		err := notifyMe(notifyAddr)
		if err != nil {
			log.Printf("Notifying error: %s", err.Error())
		} else {
			log.Printf("Successful notifying")
		}
	})

	registerReceiver(application.Calendar, notifyAddr, expireAfter)
	err := notifyMe(notifyAddr)
	if err != nil {
		log.Printf("Notifying at init failed %s", err.Error())
	}
}

func registerReceiver(cal interfaces.CalendarInterface, notifyAddr string, expireAfter int) {
	//  refresh after every some constant time interval? -- done I think
	// also there are some refreshing tokens flying around soo... yeeeah...

	err := cal.WatchForChanges(notifyAddr, time.Duration(expireAfter)*time.Second)
	if err != nil {
		log.Printf("Error sending watch request: %s", err.Error())

		go func() {
			log.Printf("Retrying in one minute")
			timer := time.NewTimer(time.Duration(time.Minute))
			refreshTime := <- timer.C
			log.Printf("Refreshing watch channel on %s", refreshTime.Format(utils.DefaultTimeType))
			registerReceiver(cal, notifyAddr, expireAfter)
		}()
		return
	}

	// if everything went smoothly, carry on with usual refresh-after-an-hour-or-so
	go func() {
		log.Printf("Scheduled for retrying in %d second", expireAfter)
		timer := time.NewTimer(time.Duration(expireAfter)*time.Second)
		refreshTime := <- timer.C
		log.Printf("Refreshing watch channel on %s", refreshTime.Format(utils.DefaultTimeType))
		registerReceiver(cal, notifyAddr, expireAfter)
	}()
}

func notifyMe(notifyAddr string) (error) {
	_, err := http.DefaultClient.Get(notifyAddr)
	return err
}
