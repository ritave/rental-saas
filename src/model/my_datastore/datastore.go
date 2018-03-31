package my_datastore

import (
	"rental-saas/src/model"
	"rental-saas/src/utils/config"
	"database/sql"
		_ "github.com/mattn/go-sqlite3"
	"errors"
	"github.com/sirupsen/logrus"
)

// TODO create table on first run

const (
	dbFile             = "calendar.db"
	sqlTableEventsJSON = `
		DROP TABLE events;
		CREATE TABLE events (
			uuid INTEGER NOT NULL PRIMARY KEY,
			jsonifiedObject TEXT
		);
	`
	sqlCreateTableEvents = `
		CREATE TABLE events (
			uuid TEXT NOT NULL PRIMARY KEY,
			user TEXT,
			start_date TEXT,
			end_date TEXT,
			creation_date TEXT,
			summary TEXT,
			location TEXT,
			timestamp_ms INT
		);
	`
	sqlDropTableEvents = `
		DROP TABLE events;
	`
	sqlQueryAll = `
		SELECT * FROM events;
	`
	sqlCountAll = `
		SELECT count(*) FROM events;
	`
	sqlDeleteEvent = `
		DELETE FROM events WHERE uuid = ?;
	`
	sqlInsertEvent = `INSERT INTO events
		(uuid, user, start_date, end_date, creation_date, summary, location, timestamp_ms) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?);
	`
)

// (uuid, user, start_date, end_date, creation_date, summary, location, timestamp_ms)

type Datastore struct {
	dbFile string
	db     *sql.DB
}

func (ds *Datastore) SynchroniseDatastore([]*model.EventModified) (SynchEffect) {
	panic("implement me")
}

func (ds *Datastore) QueryEvents() ([]*model.Event, error) {
	rows, err := ds.db.Query(sqlQueryAll)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	var count int
	countRow, err := ds.db.Query(sqlCountAll)
	if err != nil {
		logrus.Printf("Oh ffs: %s", err.Error())
		return nil, err
	}

	countRow.Next()
	countRow.Scan(&count)
	if count == 0 {
		return nil, errors.New("empty")
	}

	result := make([]*model.Event, count)
	i := 0
	for rows.Next() {
		ev, err := RowToEvent(rows)
		if err != nil {
			// TODO what do?
			logrus.Printf("Extracting failed: %s", err.Error())
			continue
		}
		result[i] = ev
		i ++
	}

	return result, nil
}

func (ds *Datastore) DeleteEvent(UUID string) (error) {
	_, err := ds.db.Exec(sqlDeleteEvent, UUID)
	return err
}

func (ds *Datastore) SaveEvent(event *model.Event) (error) {
	_, err := ds.db.Exec(sqlInsertEvent,
		event.UUID,
		event.User,
		event.Start,
		event.End,
		event.CreationDate,
		event.Summary,
		event.Location,
		event.Timestamp,
	)
	return err
}

func RowToEvent(rows *sql.Rows) (*model.Event, error) {
	var r model.Event

	rows.Next()
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

	return &r, nil
}

func (ds *Datastore) GetEvent(UUID string) (*model.Event, error) {
	getFirst := `SELECT * FROM events WHERE uuid = ?`
	rows, err := ds.db.Query(getFirst, UUID)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	return RowToEvent(rows)
}

func (ds *Datastore) Restart() {
	ds.db.Exec(sqlDropTableEvents)
	ds.db.Exec(sqlCreateTableEvents)
}

func New(c config.C) *Datastore {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		logrus.Fatal(err)
	}
	//defer db.Close()

	return &Datastore{
		db:     db,
		dbFile: dbFile,
	}
}
