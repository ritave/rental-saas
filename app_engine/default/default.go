package main

import (
	"net/http"
	"google.golang.org/appengine"
	"log"
	"rental-saas/src/view/calendar"
	"rental-saas/src/view/calendar/event"
	"github.com/rs/cors"
)

func main() {
	bindEndpoints()
	appengine.Main()
}

func bindEndpoints() {
	mux := http.NewServeMux()

	// events related
	mux.HandleFunc("/calendar/event/create", event.Create)
	mux.HandleFunc("/calendar/event/delete", event.Delete)

	// calendar related
	mux.HandleFunc("/calendar/changed", calendar.Changed)
	mux.HandleFunc("/calendar/view", calendar.View)

	// cors
	c := cors.New(cors.Options{
		AllowedOrigins: []string{CORSnpmdev, CORSappengine, CORSdeployed},
	})
	handler := c.Handler(mux)
	http.Handle("/", handler)

	log.Println("Bound endpoints.")
}

const (
	CORSnpmdev    = "http://localhost:5000"
	CORSappengine = "http://localhost:8080"
	CORSdeployed  = "https://calendarcron.appspot.com"
)
