package main

import (
	"html/template"
	"net/http"
	"time"

	"golang.org/x/net/context"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"encoding/json"
	stdLog "log"
	"calendar-synch/handlers"
	"calendar-synch/helpers"
)

var lastKey *datastore.Key

type WhatWeReallyWantToStoreIs struct {
	handlers.EventModification
	Received string
}

const keyKind = "EventModification"

func main() {
	http.HandleFunc("/dummy", handleMainPage)
	http.HandleFunc("/dummy/send", handleSend)
	stdLog.Println("Starting application")
	appengine.Main()
}

// eventModificationKey returns the key used for all entries.
func eventModificationKey(ctx context.Context) *datastore.Key {
	// The string "default_guestbook" here could be varied to have multiple guestbooks.
	return datastore.NewKey(ctx, keyKind, "default_event_modification", 0, nil)
}

var tpl = template.Must(template.ParseGlob("templates/*.html"))

func handleMainPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "GET requests only", http.StatusMethodNotAllowed)
		return
	}

	ctx := appengine.NewContext(r)
	tic := time.Now()
	q := datastore.NewQuery(keyKind).Limit(100)
	var eventsModified []*WhatWeReallyWantToStoreIs
	if _, err := q.GetAll(ctx, &eventsModified); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Errorf(ctx, "GetAll: %v", err)
		return
	}
	log.Infof(ctx, "Datastore lookup took %s", time.Since(tic).String())
	log.Infof(ctx, "Rendering %d modifications", len(eventsModified))

	if err := tpl.ExecuteTemplate(w, "dummy.html", eventsModified); err != nil {
		log.Errorf(ctx, "%v", err)
	}
}

func handleSend(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "POST requests only", http.StatusMethodNotAllowed)
		return
	}
	ctx := appengine.NewContext(r)

	eventsChanged, err := extractEventsFromBody(r)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	for _, eventChanged := range eventsChanged {
		toStore := WhatWeReallyWantToStoreIs{
			Received: helpers.TimeToString(time.Now()),
			EventModification: eventChanged,
		}

		key := datastore.NewIncompleteKey(ctx, keyKind, nil)
		if lastPut, err := datastore.Put(ctx, key, &toStore); err != nil {
			log.Errorf(ctx, err.Error())
			continue
		} else {
			lastKey = lastPut
		}
	}

}

func extractEventsFromBody(r *http.Request) (handlers.EventChangedResponse, error) {
	var target = make(handlers.EventChangedResponse, 0)
	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&target)
	if err != nil {
		return nil, err
	}
	return target, nil
}
