package tests

import (
	"reflect"
	"testing"
	"calendar-synch/logic"
	"calendar-synch/objects"
	"time"
	"calendar-synch/helpers"
	"log"
)

var zeroth = time.Now()
var first = zeroth.Add(time.Hour)
var second = first.Add(time.Hour)
var third = second.Add(time.Hour)
var fourth = third.Add(time.Hour)
var fifth = fourth.Add(time.Hour)
var sixth = fifth.Add(time.Hour)

func TestCompareSortable(t *testing.T) {

	tts := helpers.TimeToString
	
	var timeBack []time.Time
	var tBS []string
	const slots = 8
	timeBack = make([]time.Time, slots)
	tBS = make([]string, slots)
	startInPast := time.Now().Add(-slots * time.Hour)
	for i:=0; i<slots; i++ {
		timeBack[i] = startInPast.Add(time.Hour)
		tBS[i] = tts(timeBack[i])
	}

	var exhibit1 = &objects.Event{"summary", "user1@mail.com", tts(zeroth), tts(first), "location1", tBS[0]}
	var exhibit1ModifiedTimeForward = &objects.Event{"summary", "user1@mail.com", tts(first), tts(second), "location1", tBS[0]}
	var exhibit1ModifiedTimeForwardAndPlace = &objects.Event{"summary", "user1@mail.com", tts(first), tts(second), "location1-modified", tBS[0]}
	log.Printf("exhibit 1 %p %v\n", exhibit1, *exhibit1)

	var exhibit2 = &objects.Event{"summary", "user2@mail.com", tts(first), tts(second), "location2", tBS[1]}
	var exhibit2ModifiedTimeBackward = &objects.Event{"summary", "user2@mail.com", tts(zeroth), tts(first), "location2", tBS[1]}
	var exhibit2ModifiedTimeBackwardAndPlace = &objects.Event{"summary", "user2@mail.com", tts(zeroth), tts(first), "location2-modified", tBS[1]}
	log.Printf("exhibit 2 %p %v\n", exhibit2, *exhibit2)

	var exhibit3 = &objects.Event{"summary", "user3@mail.com", tts(second), tts(third), "location3", tBS[2]}
	var exhibit3ModifiedPlace = &objects.Event{"summary", "user3@mail.com", tts(second), tts(third), "location3-modified", tBS[2]}
	log.Printf("exhibit 3 %p %v\n", exhibit3, *exhibit3)

	var exhibit4 = &objects.Event{"summary", "user4@mail.com", tts(third), tts(fourth), "location4", tBS[3]}
	var exhibit4SecondEvent = &objects.Event{"summary", "user4@mail.com", tts(fifth), tts(sixth), "location4-some-other", tBS[6]}
	var exhibit4ThirdEvent = &objects.Event{"summary", "user4@mail.com", tts(fifth), tts(sixth), "location4", tBS[7]}
	log.Printf("exhibit 4_1 %p %v\n", exhibit4, *exhibit4)
	log.Printf("exhibit 4_2 %p %v\n", exhibit4SecondEvent, *exhibit4SecondEvent)
	log.Printf("exhibit 4_3 %p %v\n", exhibit4ThirdEvent, *exhibit4ThirdEvent)

	var exhibit5 = &objects.Event{"summary", "user5@mail.com", tts(fourth), tts(fifth), "location5", tBS[4]}

	var exhibit6 = &objects.Event{"summary", "user6@mail.com", tts(fifth), tts(sixth), "location6", tBS[5]}

	type args struct {
		saved  objects.SortableEvents
		actual objects.SortableEvents
	}
	tests := []struct {
		name    string
		args    args
		want    []*objects.EventModified
		wantErr bool
	}{
		{"empty", args{objects.SortableEvents{}, objects.SortableEvents{}}, []*objects.EventModified{}, false},

		{"nothing changed", args{
			objects.SortableEvents{exhibit1, exhibit2, exhibit3, exhibit4, exhibit5, exhibit6},
			objects.SortableEvents{exhibit1, exhibit2, exhibit3, exhibit4, exhibit5, exhibit6},
		}, []*objects.EventModified{}, false},

		{"2 added inside", args{
			objects.SortableEvents{exhibit1, exhibit3, exhibit4, exhibit5, exhibit6},
			objects.SortableEvents{exhibit1, exhibit2, exhibit3, exhibit4, exhibit5, exhibit6},
		}, []*objects.EventModified{objects.NewModified(exhibit2).Flag(objects.Added)}, false},

		{"2 deleted inside", args{
			objects.SortableEvents{exhibit1, exhibit2, exhibit3, exhibit4, exhibit5, exhibit6},
			objects.SortableEvents{exhibit1, exhibit3, exhibit4, exhibit5, exhibit6},
		}, []*objects.EventModified{objects.NewModified(exhibit2).Flag(objects.Deleted)}, false},

		{"1 modified time forwards", args{
			objects.SortableEvents{exhibit1, exhibit2, exhibit3, exhibit4, exhibit5, exhibit6},
			objects.SortableEvents{exhibit1ModifiedTimeForward, exhibit2, exhibit3, exhibit4, exhibit5, exhibit6},
		}, []*objects.EventModified{objects.NewModified(exhibit1ModifiedTimeForward).Flag(objects.ModifiedTime)}, false},

		{"2 modified time backwards", args{
			objects.SortableEvents{exhibit1, exhibit2, exhibit3, exhibit4, exhibit5, exhibit6},
			objects.SortableEvents{exhibit1, exhibit2ModifiedTimeBackward, exhibit3, exhibit4, exhibit5, exhibit6},
		}, []*objects.EventModified{objects.NewModified(exhibit2ModifiedTimeBackward).Flag(objects.ModifiedTime)}, false},

		{"1 modified time forwards and place", args{
			objects.SortableEvents{exhibit1, exhibit2, exhibit3, exhibit4, exhibit5, exhibit6},
			objects.SortableEvents{exhibit1ModifiedTimeForwardAndPlace, exhibit2, exhibit3, exhibit4, exhibit5, exhibit6},
		}, []*objects.EventModified{objects.NewModified(exhibit1ModifiedTimeForwardAndPlace).Flag(objects.ModifiedTime).Flag(objects.ModifiedLocation)}, false},

		{"2 modified time backwards and place", args{
			objects.SortableEvents{exhibit1, exhibit2, exhibit3, exhibit4, exhibit5, exhibit6},
			objects.SortableEvents{exhibit1, exhibit2ModifiedTimeBackwardAndPlace, exhibit3, exhibit4, exhibit5, exhibit6},
		}, []*objects.EventModified{objects.NewModified(exhibit2ModifiedTimeBackwardAndPlace).Flag(objects.ModifiedTime).Flag(objects.ModifiedLocation)}, false},

		{"3 modified place", args{
			objects.SortableEvents{exhibit1, exhibit2, exhibit3, exhibit4, exhibit5, exhibit6},
			objects.SortableEvents{exhibit1, exhibit2, exhibit3ModifiedPlace, exhibit4, exhibit5, exhibit6},
		}, []*objects.EventModified{objects.NewModified(exhibit3ModifiedPlace).Flag(objects.ModifiedLocation)}, false},

		{"5 added, 3 deleted and 2 modified", args{
			objects.SortableEvents{exhibit1, exhibit2, exhibit3, exhibit4, exhibit6},
			objects.SortableEvents{exhibit1, exhibit2ModifiedTimeBackward, exhibit4, exhibit5, exhibit6},
		}, []*objects.EventModified{
			objects.NewModified(exhibit2ModifiedTimeBackward).Flag(objects.ModifiedTime),
			objects.NewModified(exhibit3).Flag(objects.Deleted),
			objects.NewModified(exhibit5).Flag(objects.Added),
		}, false},

		{"1 and 2 swapped places", args{
			objects.SortableEvents{exhibit1, exhibit2, exhibit3, exhibit4, exhibit5, exhibit6},
			objects.SortableEvents{exhibit2ModifiedTimeBackward, exhibit1ModifiedTimeForward, exhibit3, exhibit4, exhibit5, exhibit6},
		}, []*objects.EventModified{
			objects.NewModified(exhibit2ModifiedTimeBackward).Flag(objects.ModifiedTime),
			objects.NewModified(exhibit1ModifiedTimeForward).Flag(objects.ModifiedTime),
		}, false},

		{"4 added two more events", args{
			objects.SortableEvents{exhibit1, exhibit2, exhibit3, exhibit4, exhibit5, exhibit6},
			objects.SortableEvents{exhibit1, exhibit2, exhibit3, exhibit4, exhibit5, exhibit6, exhibit4SecondEvent, exhibit4ThirdEvent},
		}, []*objects.EventModified{
			objects.NewModified(exhibit4SecondEvent).Flag(objects.Added),
			objects.NewModified(exhibit4ThirdEvent).Flag(objects.Added),
		}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := logic.CompareSortable(tt.args.saved, tt.args.actual)
			if (err != nil) != tt.wantErr {
				t.Errorf("CompareSortable() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CompareSortable() = %v, want %v", got, tt.want)
			}
		})
	}
}
