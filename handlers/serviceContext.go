package handlers

import (
	"google.golang.org/api/calendar/v3"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"net/http"
	"google.golang.org/appengine"
	"log"
)

const secretsLocation = "../secrets"

func GetService(r *http.Request) *calendar.Service {
	if !appengine.IsDevAppServer() {
		client, err := google.DefaultClient(appengine.NewContext(r), calendar.CalendarScope)
		if err != nil {
			log.Fatalf("Default client failed: %s", err.Error())
		}
		srv, err := calendar.New(client)
		if err != nil {
			log.Fatalf("Creating calendar service on the spot failed: %s", err.Error())
		}
		return srv
	} else {
		b, err := ioutil.ReadFile(secretsLocation + "/service_client_default.json")
		if err != nil {
			log.Fatalf("Unable to read client secret file: %v", err)
		}

		config, err := google.JWTConfigFromJSON(b, calendar.CalendarScope)
		if err != nil {
			log.Fatalf("Unable to parse service client secret file to config: %v", err)
		}
		client := config.Client(appengine.NewContext(r))

		srv, err := calendar.New(client)
		if err != nil {
			log.Fatalf("Unable to retrieve calendar Client %v", err)
		}
		return srv
	}
}
