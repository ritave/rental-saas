package main

import (
	"net/http"
	"google.golang.org/appengine"
	"log"
	"calendar-synch/handlers"
	"html/template"
)

func main() {
	bindEndpoints()
	appengine.Main()
}

func bindEndpoints() {
	http.HandleFunc("/event/ping", Ping)
	http.HandleFunc("/event/create", handlers.EventCreate)
	http.HandleFunc("/event/list", handlers.EventList)
	http.HandleFunc("/event/changed", handlers.EventChanged)

	http.HandleFunc("/calendar/view", handlers.CalendarView)

	log.Println("Bound endpoints...")
}

func init() {
	// yes, I know it's a no-op
	log.Println("Initialized stuff...")
}

var tempTemplate = template.Must(template.New("temp").Parse(`
<html>
  <head>
    <title>Hey, hello, welcome</title>
  </head>
  <body>
    <pre>This should be on "/event/"</pre>
	<p>From {{.From}}</p>
	<p>To {{.To}}</p>
  </body>
</html>
`))
type RequestData struct {
	From      string   `json:"from"`
	To        string   `json:"to"`
}
func Ping(w http.ResponseWriter, r *http.Request) {
	requestData := &RequestData{From: r.RemoteAddr, To: r.Host}

	if err := tempTemplate.Execute(w, requestData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}