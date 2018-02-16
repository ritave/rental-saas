package event

import (
	"net/http"
	"google.golang.org/appengine"
	"log"
	"calendar-synch/handlers"
)

// [START notification_struct]

func main() {
	bindEndpoints()
	appengine.Main()
}

func bindEndpoints() {
	http.HandleFunc("/", handlers.Root)
	http.HandleFunc("/notify", handlers.Notify)
	http.HandleFunc("/event/create", handlers.EventCreate)
	log.Println("Bound endpoints...")
}

func init() {
	log.Println("Initialized stuff...")
}
