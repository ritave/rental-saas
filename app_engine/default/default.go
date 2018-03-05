package main

import (
	"net/http"
	"google.golang.org/appengine"
	"log"
	"html/template"
	"calendar-synch/src/handlers/calendar"
	"calendar-synch/src/handlers/event"
)

func main() {
	bindEndpoints()
	appengine.Main()
}

func bindEndpoints() {
	http.HandleFunc("/event/ping", Ping) // TODO delete/move/update/idk
	http.HandleFunc("/event/create", event.Create)
	http.HandleFunc("/event/delete", event.Delete)
	http.HandleFunc("/event/changed", event.Changed)

	http.HandleFunc("/event/list", event.List) // TODO move to calendar/view or /list or /get idc
	http.HandleFunc("/calendar/view", calendar.View)

	log.Println("Bound endpoints...")
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