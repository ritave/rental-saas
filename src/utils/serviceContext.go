package utils

import (
	"google.golang.org/api/calendar/v3"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"net/http"
	"google.golang.org/appengine"
		"context"
	"github.com/sirupsen/logrus"
)

const secretsLocation = "secrets"

func NewStandard(r *http.Request) *calendar.Service {
	//if !appengine.IsDevAppServer() {
	//	client, err := google.DefaultClient(appengine.NewContext(r), calendar.CalendarScope)
	//	if err != nil {
	//		logrus.Fatalf("Default client failed: %s", err.Error())
	//	}
	//	srv, err := calendar.New(client)
	//	if err != nil {
	//		logrus.Fatalf("Creating calendar service on the spot failed: %s", err.Error())
	//	}
	//	return srv
	//} else {
	b, err := ioutil.ReadFile(secretsLocation + "/service_client.json")
	if err != nil {
		logrus.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.JWTConfigFromJSON(b, calendar.CalendarScope)
	if err != nil {
		logrus.Fatalf("Unable to parse service client secret file to config: %v", err)
	}
	client := config.Client(appengine.NewContext(r))

	srv, err := calendar.New(client)
	if err != nil {
		logrus.Fatalf("Unable to retrieve calendar Client %v", err)
	}
	return srv
	//}
}

func NewFlex(ctx context.Context) *calendar.Service {
	b, err := ioutil.ReadFile(secretsLocation + "/service_client.json")
	if err != nil {
		logrus.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.JWTConfigFromJSON(b, calendar.CalendarScope)
	if err != nil {
		logrus.Fatalf("Unable to parse service client secret file to config: %v", err)
	}
	client := config.Client(ctx)

	srv, err := calendar.New(client)
	if err != nil {
		logrus.Fatalf("Unable to retrieve calendar Client %v", err)
	}
	return srv
}
