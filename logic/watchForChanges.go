package logic

import (
	"google.golang.org/api/calendar/v3"
	"github.com/satori/go.uuid"
	"time"
	"log"
)

type ImportantChannelFields struct {
	ResourceId string
	Uuid string
}

func WatchForChanges(cal *calendar.Service, receiver string, expireAfter time.Duration) (error) {
	watchChannel, err := newChannel(cal, receiver, expireAfter)

	receipt := ImportantChannelFields{
		ResourceId: watchChannel.ResourceId,
		Uuid: watchChannel.Id,
	}

	log.Printf("ResourceId: %s | Id: %s ", receipt.ResourceId, receipt.Uuid)

	return err
}

func stopChannel(cal *calendar.Service, resourceID, uuid string) (error) {
	return cal.Channels.Stop(
		&calendar.Channel{
			ResourceId: resourceID,
			Id: uuid,
		},
	).Do()
}

func newChannel(cal *calendar.Service, receiver string, expireAfter time.Duration) (*calendar.Channel, error) {
	u := uuid.Must(uuid.NewV4())

	channel := calendar.Channel{
		Id: u.String(),
		Address: receiver,
		Type: "web_hook",
		Expiration: time.Now().Add(expireAfter).UnixNano(),
	}
	return cal.Events.Watch("primary", &channel).Do()
}