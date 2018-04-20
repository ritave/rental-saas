package tests

import (
	"reflect"
	"testing"
	"rental-saas/service/src/presenter"
	"rental-saas/service/src/model"
	"time"
	"rental-saas/service/src/utils"
	"github.com/satori/go.uuid"
)

var zeroth = time.Now()
var first = zeroth.Add(time.Hour)
var second = first.Add(time.Hour)
var third = second.Add(time.Hour)
var fourth = third.Add(time.Hour)
var fifth = fourth.Add(time.Hour)
var sixth = fifth.Add(time.Hour)

func TestCompareSortable(t *testing.T) {

	tts := utils.TimeToString
	uuidString := uuid.Must(uuid.NewV4()).String
	
	var timeBack []time.Time
	var tBS []string
	const slots = 8
	timeBack = make([]time.Time, slots)
	tBS = make([]string, slots)
	startInPast := time.Now().Add(-slots * time.Hour)
	for i:=0; i<slots; i++ {
		timeBack[i] = startInPast.Add(time.Duration(i) * time.Hour)
		tBS[i] = tts(timeBack[i])
	}

	// YO DAWG I HERD YOU LIKE VARIABLES SO WE PUT VARIABLES IN YO VARIABLES SO YOU CAN COMPUTE WHILE YOU COMPUTE

	var xzibit1 = &model.Event{"summary", "user1@mail.com", tts(zeroth), tts(first), "location1", tBS[0], 0, uuidString(), ""}
	var xzibit1ModifiedTimeForward = &model.Event{"summary", "user1@mail.com", tts(first), tts(second), "location1", tBS[0], 0, uuidString(), ""}
	var xzibit1ModifiedTimeForwardAndPlace = &model.Event{"summary", "user1@mail.com", tts(first), tts(second), "location1-modified", tBS[0], 0, uuidString(), ""}

	var xzibit2 = &model.Event{"summary", "user2@mail.com", tts(first), tts(second), "location2", tBS[1], 1, uuidString(), ""}
	var xzibit2ModifiedTimeBackward = &model.Event{"summary", "user2@mail.com", tts(zeroth), tts(first), "location2", tBS[1], 1, uuidString(), ""}
	var xzibit2ModifiedTimeBackwardAndPlace = &model.Event{"summary", "user2@mail.com", tts(zeroth), tts(first), "location2-modified", tBS[1], 1, uuidString(), ""}

	var xzibit3 = &model.Event{"summary", "user3@mail.com", tts(second), tts(third), "location3", tBS[2], 2, uuidString(), ""}
	var xzibit3ModifiedPlace = &model.Event{"summary", "user3@mail.com", tts(second), tts(third), "location3-modified", tBS[2], 2, uuidString(), ""}

	var xzibit4 = &model.Event{"summary", "user4@mail.com", tts(third), tts(fourth), "location4", tBS[3], 3, uuidString(), ""}
	var xzibit4SecondEvent = &model.Event{"summary", "user4@mail.com", tts(fifth), tts(sixth), "location4-some-other", tBS[6], 6, uuidString(), ""}
	var xzibit4ThirdEvent = &model.Event{"summary", "user4@mail.com", tts(fifth), tts(sixth), "location4", tBS[7], 7, uuidString(), ""}

	var xzibit5 = &model.Event{"summary", "user5@mail.com", tts(fourth), tts(fifth), "location5", tBS[4], 4, uuidString(), ""}

	var xzibit6 = &model.Event{"summary", "user6@mail.com", tts(fifth), tts(sixth), "location6", tBS[5], 5, uuidString(), ""}

	//log.Printf("xzibit 1 %p %v\n", xzibit1, *xzibit1)
	//log.Printf("xzibit 2 %p %v\n", xzibit2, *xzibit2)
	//log.Printf("xzibit 3 %p %v\n", xzibit3, *xzibit3)
	//log.Printf("xzibit 4_1 %p %v\n", xzibit4, *xzibit4)
	//log.Printf("xzibit 4_2 %p %v\n", xzibit4SecondEvent, *xzibit4SecondEvent)
	//log.Printf("xzibit 4_3 %p %v\n", xzibit4ThirdEvent, *xzibit4ThirdEvent)


	type args struct {
		saved  model.SortableEvents
		actual model.SortableEvents
	}
	tests := []struct {
		name    string
		args    args
		want    []*model.EventModified
		wantErr bool
	}{
		{"empty", args{model.SortableEvents{}, model.SortableEvents{}}, []*model.EventModified{}, false},

		{"nothing changed", args{
			model.SortableEvents{xzibit1, xzibit2, xzibit3, xzibit4, xzibit5, xzibit6},
			model.SortableEvents{xzibit1, xzibit2, xzibit3, xzibit4, xzibit5, xzibit6},
		}, []*model.EventModified{}, false},

		{"2 added inside", args{
			model.SortableEvents{xzibit1, xzibit3, xzibit4, xzibit5, xzibit6},
			model.SortableEvents{xzibit1, xzibit2, xzibit3, xzibit4, xzibit5, xzibit6},
		}, []*model.EventModified{model.NewModified(xzibit2).Flag(model.Added)}, false},

		{"2 deleted inside", args{
			model.SortableEvents{xzibit1, xzibit2, xzibit3, xzibit4, xzibit5, xzibit6},
			model.SortableEvents{xzibit1, xzibit3, xzibit4, xzibit5, xzibit6},
		}, []*model.EventModified{model.NewModified(xzibit2).Flag(model.Deleted)}, false},

		{"1 modified time forwards", args{
			model.SortableEvents{xzibit1, xzibit2, xzibit3, xzibit4, xzibit5, xzibit6},
			model.SortableEvents{xzibit1ModifiedTimeForward, xzibit2, xzibit3, xzibit4, xzibit5, xzibit6},
		}, []*model.EventModified{model.NewModified(xzibit1ModifiedTimeForward).Flag(model.ModifiedTime)}, false},

		{"2 modified time backwards", args{
			model.SortableEvents{xzibit1, xzibit2, xzibit3, xzibit4, xzibit5, xzibit6},
			model.SortableEvents{xzibit1, xzibit2ModifiedTimeBackward, xzibit3, xzibit4, xzibit5, xzibit6},
		}, []*model.EventModified{model.NewModified(xzibit2ModifiedTimeBackward).Flag(model.ModifiedTime)}, false},

		{"1 modified time forwards and place", args{
			model.SortableEvents{xzibit1, xzibit2, xzibit3, xzibit4, xzibit5, xzibit6},
			model.SortableEvents{xzibit1ModifiedTimeForwardAndPlace, xzibit2, xzibit3, xzibit4, xzibit5, xzibit6},
		}, []*model.EventModified{model.NewModified(xzibit1ModifiedTimeForwardAndPlace).Flag(model.ModifiedTime).Flag(model.ModifiedLocation)}, false},

		{"2 modified time backwards and place", args{
			model.SortableEvents{xzibit1, xzibit2, xzibit3, xzibit4, xzibit5, xzibit6},
			model.SortableEvents{xzibit1, xzibit2ModifiedTimeBackwardAndPlace, xzibit3, xzibit4, xzibit5, xzibit6},
		}, []*model.EventModified{model.NewModified(xzibit2ModifiedTimeBackwardAndPlace).Flag(model.ModifiedTime).Flag(model.ModifiedLocation)}, false},

		{"3 modified place", args{
			model.SortableEvents{xzibit1, xzibit2, xzibit3, xzibit4, xzibit5, xzibit6},
			model.SortableEvents{xzibit1, xzibit2, xzibit3ModifiedPlace, xzibit4, xzibit5, xzibit6},
		}, []*model.EventModified{model.NewModified(xzibit3ModifiedPlace).Flag(model.ModifiedLocation)}, false},

		{"5 added, 3 deleted and 2 modified", args{
			model.SortableEvents{xzibit1, xzibit2, xzibit3, xzibit4, xzibit6},
			model.SortableEvents{xzibit1, xzibit2ModifiedTimeBackward, xzibit4, xzibit5, xzibit6},
		}, []*model.EventModified{
			model.NewModified(xzibit2ModifiedTimeBackward).Flag(model.ModifiedTime),
			model.NewModified(xzibit3).Flag(model.Deleted),
			model.NewModified(xzibit5).Flag(model.Added),
		}, false},

		{"1 and 2 swapped places", args{
			model.SortableEvents{xzibit1, xzibit2, xzibit3, xzibit4, xzibit5, xzibit6},
			model.SortableEvents{xzibit2ModifiedTimeBackward, xzibit1ModifiedTimeForward, xzibit3, xzibit4, xzibit5, xzibit6},
		}, []*model.EventModified{
			model.NewModified(xzibit1ModifiedTimeForward).Flag(model.ModifiedTime),
			model.NewModified(xzibit2ModifiedTimeBackward).Flag(model.ModifiedTime),
		}, false},

		{"4 added two more events", args{
			model.SortableEvents{xzibit1, xzibit2, xzibit3, xzibit4, xzibit5, xzibit6},
			model.SortableEvents{xzibit1, xzibit2, xzibit3, xzibit4, xzibit5, xzibit6, xzibit4SecondEvent, xzibit4ThirdEvent},
		}, []*model.EventModified{
			model.NewModified(xzibit4SecondEvent).Flag(model.Added),
			model.NewModified(xzibit4ThirdEvent).Flag(model.Added),
		}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := presenter.CompareSortable(tt.args.saved, tt.args.actual)
			if (err != nil) != tt.wantErr {
				t.Errorf("CompareSortable() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CompareSortable() =\n")
				if len(got) == 0 {
					t.Errorf("[]\n")
				} else {
					//t.Errorf("[\n")
					for _, el := range got {
						t.Errorf("E: %v, M: %v\n", el.Event, el.Modifications)
					}
					//t.Errorf("]\n")
				}

				t.Errorf("want\n")
				if len(tt.want) == 0 {
					t.Errorf("[]\n")
				} else {
					//t.Errorf("[\n")
					for _, el := range tt.want {
						t.Errorf("E: %v, M: %v\n", el.Event, el.Modifications)
					}
					//t.Errorf("]\n")
				}
			}
		})
	}
}
