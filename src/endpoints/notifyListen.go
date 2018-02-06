package endpoints

import (
	"io/ioutil"
	"fmt"
	"time"
	"net/http"
	"google.golang.org/appengine"
	"google.golang.org/appengine/user"
	"google.golang.org/appengine/datastore"
	"context"
	"google.golang.org/api/calendar/v3"
)

type NotificationEntity struct {
	Source  string
	Content []byte
	Date    time.Time
}

func NotifyListen(s *calendar.Service, w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("Error parsing response: %S\n", err)
	}

	fmt.Printf("%S\n", string(body))

	g := NotificationEntity{
		Content: body,
		Date:    time.Now(),
	}
	if u := user.Current(c); u != nil {
		g.Source = u.String()
	}
	// We set the same parent key on every NotificationEntity entity to ensure each NotificationEntity
	// is in the same entity group. Queries across the single entity group
	// will be consistent. However, the write rate to a single entity group
	// should be limited to ~1/second.
	key := datastore.NewIncompleteKey(c, "NotificationEntity", logkKey(c))
	_, err = datastore.Put(c, key, &g)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

// logkKey returns the key used for all notification entries.
func logkKey(c context.Context) *datastore.Key {
	// The string "default_guestbook" here could be varied to have multiple notifications.
	return datastore.NewKey(c, "NotificationEntity", "default_notification", 0, nil)
}
