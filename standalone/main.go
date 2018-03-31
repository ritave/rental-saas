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
	. "rental-saas/src/application/core"
	. "rental-saas/src/application/handler"
	"rental-saas/src/utils/config"
	"rental-saas/src/application/interfaces"
	"strconv"
	"github.com/sirupsen/logrus"
	"rental-saas/src/view"
)

func main() {
	app := New(config.GetConfig())
	mux := http.NewServeMux()

	// logging init
	if app.Config.Logging.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	// events related
	mux.Handle("/calendar/event/create", &AppHandler{app, event.CreateRequest{}, event.Create})
	mux.Handle("/calendar/event/delete", &AppHandler{app, event.DeleteRequest{}, event.Delete})

	// calendar related
	mux.Handle("/calendar/changed", &AppHandler{app, calendar.ChangedRequest{}, calendar.Changed})
	mux.Handle("/calendar/view", &AppHandler{app, calendar.ViewRequest{}, calendar.View})
	
	// notify related
	mux.Handle("/notify/get", &AppHandler{app, notify.GetRequest{}, notify.Get})
	mux.HandleFunc("/notify/channel/delete", notify.DeleteChannel)

	// keep alive & admin retarted
	mux.HandleFunc("/ping", view.HandlerPing)

	// cors
	c := cors.New(cors.Options{
		AllowedOrigins: app.Config.CORS,
	})
	handler := c.Handler(mux)
	http.Handle("/", handler)

	// datastore
	if app.Config.DB.Restart {
		log.Println("Dropping and creating tables in database")
		app.Datastore.Restart()
	}

	// watch notify
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

	logrus.Printf("Listening on port %d", app.Config.Server.Port)
	logrus.Fatal(http.ListenAndServe(":"+strconv.Itoa(app.Config.Server.Port), nil))
}


// FIXME
// some ugly fuckery going on down here
func notifySetup(application *Application) {
	notifyAddr := application.Config.Receiver.Channel
	expireAfter := application.Config.Receiver.Expiration
	const threeSoundsFine = 3

	// this thing aggregates multiple requests that may come at once
	application.Utils.Ticker = utils.New(threeSoundsFine*time.Second, func(){
		calendar.Changed(application, calendar.ChangedRequest{})
	})

	registerReceiver(application.Calendar, notifyAddr, expireAfter)
	calendar.Changed(application, calendar.ChangedRequest{})
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
