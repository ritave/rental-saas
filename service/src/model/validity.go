package model

import (
	"regexp"
	"rental-saas/service/src/utils"
	"errors"
	"time"
)

var emailParser = regexp.MustCompile(`(^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$)`)

func ValidateEventFromRequest(ev Event) (*Event, error) {
	var result = Event(ev)

	startT, err := utils.VerifyStringToTime(ev.Start)
	if err != nil {
		return nil, errors.New("invalid datetime format (accepted is RFC3339: 2006-01-02T15:04:05Z or 2006-01-02T15:04:05+07:00); supplied was: "+ev.Start)
	}

	endT, err := utils.VerifyStringToTime(ev.End)
	if err != nil {
		return nil, errors.New("invalid datetime format (accepted is RFC3339: 2006-01-02T15:04:05Z or 2006-01-02T15:04:05+07:00); supplied was: "+ev.End)
	}

	if startT.Before(time.Now()) {
		return nil, errors.New("event's start cannot be set in the past")
	}

	if endT.Before(startT) {
		return nil, errors.New("event's end earlier than the beggining")
	}

	if endT.Before(time.Now()) {
		return nil, errors.New("event end cannot be set in the past")
	}

	if ev.Location == "" {
		return nil, errors.New("location cannot be empty")
	}

	if !emailParser.Match([]byte(ev.User)) {
		return nil, errors.New("invalid email address")
	}

	// FIXME temporary
	result.UserID = 1
	result.OrderID = 42

	return &result, nil
}
