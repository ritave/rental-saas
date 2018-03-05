package objects

import (
	"regexp"
	"calendar-synch/src/utils"
	"errors"
	"time"
	"log"
)

var emailParser = regexp.MustCompile(`(^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$)`)

func EvenMoreChecksForTheEvent(ev Event) (error) {
	startT, err := utils.VerifyStringToTime(ev.Start)
	if err != nil {
		log.Printf("Start date: %s", err.Error())
		return errors.New("invalid datetime format (accepted is RFC3339: 2006-01-02T15:04:05Z or 2006-01-02T15:04:05+07:00); supplied was: "+ev.Start)
	}

	endT, err := utils.VerifyStringToTime(ev.End)
	if err != nil {
		log.Printf("End date: %s", err.Error())
		return errors.New("invalid datetime format (accepted is RFC3339: 2006-01-02T15:04:05Z or 2006-01-02T15:04:05+07:00); supplied was: "+ev.End)
	}

	if startT.Before(time.Now()) {
		return errors.New("event start cannot be set in the past")
	}

	if endT.Before(startT) {
		return errors.New("event end cannot be before the start")
	}

	if endT.Before(time.Now()) {
		return errors.New("event end cannot be set in the past")
	}

	if ev.Location == "" {
		return errors.New("location cannot be empty")
	}

	if !emailParser.Match([]byte(ev.User)) {
		return errors.New("invalid email address")
	}

	return nil
}
