package model

import (
	"regexp"
	"rental-saas/src/utils"
	"errors"
	"time"
)

var emailParser = regexp.MustCompile(`(^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$)`)

func EvenMoreChecksForTheEvent(ev Event) (error) {
	startT, err := utils.VerifyStringToTime(ev.Start)
	if err != nil {
		return errors.New("invalid datetime format (accepted is RFC3339: 2006-01-02T15:04:05Z or 2006-01-02T15:04:05+07:00); supplied was: "+ev.Start)
	}

	endT, err := utils.VerifyStringToTime(ev.End)
	if err != nil {
		return errors.New("invalid datetime format (accepted is RFC3339: 2006-01-02T15:04:05Z or 2006-01-02T15:04:05+07:00); supplied was: "+ev.End)
	}

	if startT.Before(time.Now()) {
		return errors.New("event's start cannot be set in the past")
	}

	if endT.Before(startT) {
		return errors.New("event's end earlier than the beggining")
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
