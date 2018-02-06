package app_engine

import (
	"html/template"
	"net/http"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"calendar-synch/src/endpointsints"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	
	"log"
	"fmt"
)

// [START notification_struct]
type Notification struct {
	Source  string
	Content string
	Date    time.Time
}

func init() {
	fmt.Println("Where are my fucking logs")
	log.Println("I seriously need them now")
	var background = context.Background()

	account, err := appengine.ServiceAccount(context.Background())
	if err != nil {
		log.Fatalf("Service Account error: %s", err.Error())
	}

	sdkConfig, err := google.NewSDKConfig(account)
	if err != nil {
		log.Fatalf("NewSDKConfig: %s", err.Error())
	}
	service, err := calendar.New(sdkConfig.Client(background))
	if err != nil {
		log.Fatalf("New service: %s", err.Error())
	}

	http.HandleFunc("/", root)
	http.Handle("/notify", endpoints.MyHandler{service,endpoints.NotifyListen})
	http.Handle("/createEvent", endpoints.MyHandler{service, endpoints.CreateNewEvent})
}

// [START func_root]
func root(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	// Ancestor queries, as shown here, are strongly consistent with the High
	// Replication Datastore. Queries that span entity groups are eventually
	// consistent. If we omitted the .Ancestor from this query there would be
	// a slight chance that Notification that had just been written would not
	// show up in a query.
	// [START query]
	q := datastore.NewQuery("Notification").Ancestor(logkKey(c)).Order("-Date").Limit(10)
	// [END query]
	// [START getall]
	notifications := make([]Notification, 0, 10)
	if _, err := q.GetAll(c, &notifications); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// [END getall]
	if err := guestbookTemplate.Execute(w, notifications); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
// [END func_root]

var guestbookTemplate = template.Must(template.New("book").Parse(`
<html>
  <head>
    <title>Notifications store</title>
  </head>
  <body>
    {{range .}}
      {{with .Source}}
        <p><b>{{.}}</b> :</p>
      {{else}}
        <p>An anonymous person sent:</p>
      {{end}}
      <pre>{{.Date}}</pre>
      <pre>{{.Content}}</pre>
    {{end}}
  </body>
</html>
`))

// logkKey returns the key used for all guestbook entries.
func logkKey(c context.Context) *datastore.Key {
	// The string "default_guestbook" here could be varied to have multiple guestbooks.
	return datastore.NewKey(c, "Notification", "default_notification", 0, nil)
}
