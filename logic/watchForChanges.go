package logic

import (
	"google.golang.org/api/calendar/v3"
	"log"
	"github.com/satori/go.uuid"
	"time"
)

func WatchForChanges(cal *calendar.Service, receiver string) (error) {
	u := uuid.Must(uuid.NewV4())

	channel := calendar.Channel{
		Id: u.String(),
		Address: receiver,
		Type: "web_hook",
		Expiration: time.Now().Add(time.Hour).UnixNano(),
	}
	watchChannel, err := cal.Events.Watch("primary", &channel).Do()

	// TODO do something with this channel!
	if err != nil {
		log.Println("Watch response: \n", watchChannel)
		watchChannel.Expiration = time.Now().Add(time.Minute).UnixNano()
		resp, err := cal.Events.Watch("primary", watchChannel).Do()
		if err != nil {
			log.Println("Editing the channel failed")
		}
		log.Println("Second watch response: \n", resp)
	}

	// testing cancelling saved watch requests
	testChannel := calendar.Channel{
		Id: "1f672443-ef5b-4d43-a232-b36e86bb7b79",
		Type: "web_hook",
		Address: receiver,
		ResourceId: "1521728085000",
	}
	err = cal.Channels.Stop(&testChannel).Do()
	if err != nil {
		log.Println("Deleting the channel failed")
	}

	return err
}
