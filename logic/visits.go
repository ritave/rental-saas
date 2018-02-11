package logic

import (
	"time"
	"io"
	"io/ioutil"
	"golang.org/x/net/context"
	"log"
	"google.golang.org/appengine/datastore"
)

type Visit struct {
	Timestamp time.Time
	UserIP    string
	Body      []byte
}

// TODO ancestors

func RecordVisit(ctx context.Context, now time.Time, userIP string, body io.ReadCloser) error {
	defer body.Close()
	bodyContents, err := ioutil.ReadAll(body)
	if err != nil {
		log.Println("Couldn't read contents of the body:", err.Error())
	}

	v := &Visit{
		Timestamp: now,
		UserIP:    userIP,
		Body:      bodyContents,
	}

	k := datastore.NewIncompleteKey(ctx,"Visit", nil)

	_, err = datastore.Put(ctx, k, v)
	return err
}

func QueryVisits(ctx context.Context, limit int) ([]*Visit, error) {
	// Print out previous visits.
	q := datastore.NewQuery("Visit").
		Order("-Timestamp").
		Limit(limit)

	visits := make([]*Visit, 0)
	_, err := q.GetAll(ctx, &visits)
	return visits, err
}
