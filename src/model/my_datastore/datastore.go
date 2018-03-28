package my_datastore

import (
	"rental-saas/src/model"
	"rental-saas/src/utils/config"
	"database/sql"
	"log"
	"fmt"
	"errors"
	_ "github.com/mattn/go-sqlite3"
)

// TODO create table on first run

const (
	dbFile         = "calendar.db"
	sqlTableEventsJSON = `
		DROP TABLE events;
		CREATE TABLE events (
			uuid INTEGER NOT NULL PRIMARY KEY,
			jsonifiedObject TEXT
		);
		`
	sqlCreateTableEvents = `
		CREATE TABLE events (
			uuid INTEGER NOT NULL PRIMARY KEY,
			user text,
			start_date text,
			end_date text,
			creationdate text,
			summary text,
			location text,
			timestamp_ms int 
		);
		`
)

type Datastore struct {
	// not really a persistent database, lol
	ds     map[string]*model.Event
	dbFile string
	db     *sql.DB
}

func (ds *Datastore) SynchroniseDatastore([]*model.EventModified) (SynchEffect) {
	panic("implement me")
}

func (ds *Datastore) QueryEvents() ([]*model.Event, error) {
	result := make([]*model.Event, len(ds.ds))
	i := 0
	for _, v := range ds.ds {
		result[i] = v
		i ++
	}
	return result, nil
}

func (ds *Datastore) DeleteEvent(UUID string) (error) {
	ds.ds[UUID] = nil
	return nil
}

func (ds *Datastore) SaveEvent(event *model.Event) (error) {
	ds.ds[event.UUID] = event
	return nil
}

func RowsToEvent(rows *sql.Rows) (*model.Event, error) {
	var r model.Event
	if rows.Next() {
		err := rows.Scan(
			&r.UUID,
			&r.User,
			&r.Start,
			&r.End,
			&r.CreationDate,
			&r.Summary,
			&r.Location,
			&r.Timestamp,
		)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("empty rows")
	}

	return &r, nil
}

func EventToQuery(e *model.Event) (string) {
	return fmt.Sprintf("INSERT INTO events VALUES (%s, '%s', '%s', '%s', '%s', '%s', '%s', %d)",
		e.UUID,
		e.User,
		e.Start,
		e.End,
		e.CreationDate,
		e.Summary,
		e.Location,
		e.Timestamp,
	)
}

func (ds *Datastore) GetEvent(UUID string) (*model.Event, error) {
	getFirst := fmt.Sprintf(`SELECT * FROM events WHERE uuid = %s`, UUID)
	rows, err := ds.db.Query(getFirst)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	return RowsToEvent(rows)
}

func (ds *Datastore) PutEvent(event *model.Event) (error) {
	_, err := ds.db.Exec(EventToQuery(event))
	return err
}

func (ds *Datastore) dryRun() (error) {
	_, err := ds.db.Exec(sqlCreateTableEvents)
	if err != nil {
		log.Fatalf("%q: %s\n", err, sqlCreateTableEvents)
		return nil
	}

	return err
}

func New(c config.C) *Datastore {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatal(err)
	}
	//defer db.Close()


	return &Datastore{
		db:     db,
		dbFile: dbFile,
	}
}
