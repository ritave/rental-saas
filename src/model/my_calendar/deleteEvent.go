package my_calendar

import (
	"google.golang.org/api/calendar/v3"
)

func DeleteEvent(cal *calendar.Service, uuid string) (error){
	return cal.Events.Delete("primary", uuid).Do()
}
