package logic

import (
	"google.golang.org/api/calendar/v3"
	"github.com/satori/go.uuid"
	"time"
	"log"
	"calendar-synch/src/utils"
)

type ImportantChannelFields struct {
	ResourceId string
	Uuid       string
	Expiration string
	Receiver   string
}

func WatchForChanges(cal *calendar.Service, receiver string, expireAfter time.Duration) (error, ImportantChannelFields) {
	receipt := ImportantChannelFields{}

	watchChannel, err := newChannel(cal, receiver, expireAfter)
	if err != nil {
		log.Printf("New channel: %s", err.Error())
	} else {
		receipt = ImportantChannelFields{
			ResourceId: watchChannel.ResourceId,
			Uuid:       watchChannel.Id,
			Expiration: utils.TimeToString(utils.Int64ToTime(watchChannel.Expiration)),
			Receiver: watchChannel.Address,
		}

		log.Printf("ResourceId: %s | Id: %s | Receiver: %s | Expires: %s", receipt.ResourceId, receipt.Uuid, receipt.Receiver, receipt.Expiration)
	}

	return err, receipt
}

func stopChannel(cal *calendar.Service, resourceID, uuid string) (error) {
	return cal.Channels.Stop(
		&calendar.Channel{
			ResourceId: resourceID,
			Id:         uuid,
		},
	).Do()
}

func newChannel(cal *calendar.Service, receiver string, expireAfter time.Duration) (*calendar.Channel, error) {
	u := uuid.Must(uuid.NewV4())

	channel := calendar.Channel{
		Id:         u.String(),
		Address:    receiver,
		Type:       "web_hook",
		Expiration: time.Now().Add(expireAfter).UnixNano(),
	}
	return cal.Events.Watch("primary", &channel).Do()
}
