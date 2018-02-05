package notification

import (
"html/template"
"net/http"
"time"

"golang.org/x/net/context"
"google.golang.org/appengine"
"google.golang.org/appengine/datastore"
"google.golang.org/appengine/user"
	"google.golang.org/api/calendar/v3"
)

// [START notification_struct]
type Notification struct {
	Source  string
	Content string
	Date    time.Time
}

// [END notification_struct]

func init() {
	http.HandleFunc("/", root)
	http.HandleFunc("/notification", notification)
}

// guestbookKey returns the key used for all guestbook entries.
func guestbookKey(c context.Context) *datastore.Key {
	// The string "default_guestbook" here could be varied to have multiple guestbooks.
	return datastore.NewKey(c, "Guestbook", "default_guestbook", 0, nil)
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
	q := datastore.NewQuery("Notification").Ancestor(guestbookKey(c)).Order("-Date").Limit(10)
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

// [START func_sign]
func notification(w http.ResponseWriter, r *http.Request) {
	// [START new_context]
	c := appengine.NewContext(r)
	// [END new_context]



	g := Notification{
		Content: r.FormValue("content"),
		Date:    time.Now(),

	}
	// [START if_user]
	if u := user.Current(c); u != nil {
		g.Source = u.String()
	}
	// We set the same parent key on every Notification entity to ensure each Notification
	// is in the same entity group. Queries across the single entity group
	// will be consistent. However, the write rate to a single entity group
	// should be limited to ~1/second.
	key := datastore.NewIncompleteKey(c, "Notification", guestbookKey(c))
	_, err := datastore.Put(c, key, &g)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
	// [END if_user]
}

// [END func_sign]

