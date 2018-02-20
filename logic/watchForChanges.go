package logic

import (
	"google.golang.org/api/calendar/v3"
	"log"
	"math"
	"github.com/satori/go.uuid"
)

func WatchForChanges(cal *calendar.Service, receiver string) {
	u := uuid.Must(uuid.NewV4())

	channel := calendar.Channel{
		Id: u.String(),
		Address: receiver,
		Type: "web_hook",
		Expiration: math.MaxInt64, // lol
	}
	resp, err := cal.Events.Watch("primary", &channel).Do()
	if err != nil {
		log.Fatalf("Error sending watch request: %s", err.Error())
	}
	log.Println("Watch response: \n", resp)
}
