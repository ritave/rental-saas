package objects

import (
	"testing"
	"calendar-synch/src/handlers/calendar/event"
)

func TestEvent_IsTheSame(t *testing.T) {
	type fields struct {
		Summary           string
		User              string
		Start             string
		End               string
		Location          string
		CreationTimestamp int64
		CreationDate      string
		UUID              string
	}
	type args struct {
		to *Event
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Event{
				Summary:      tt.fields.Summary,
				User:         tt.fields.User,
				Start:        tt.fields.Start,
				End:          tt.fields.End,
				Location:     tt.fields.Location,
				Timestamp:    tt.fields.CreationTimestamp,
				CreationDate: tt.fields.CreationDate,
				UUID:         tt.fields.UUID,
			}
			if got := e.IsTheSame(tt.args.to); got != tt.want {
				t.Errorf("Event.IsTheSame() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEvent_Less(t *testing.T) {
	type fields struct {
		Summary           string
		User              string
		Start             string
		End               string
		Location          string
		CreationTimestamp int64
		CreationDate      string
		UUID              string
	}
	type args struct {
		than *Event
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Event{
				Summary:      tt.fields.Summary,
				User:         tt.fields.User,
				Start:        tt.fields.Start,
				End:          tt.fields.End,
				Location:     tt.fields.Location,
				Timestamp:    tt.fields.CreationTimestamp,
				CreationDate: tt.fields.CreationDate,
				UUID:         tt.fields.UUID,
			}
			if got := e.Less(tt.args.than); got != tt.want {
				t.Errorf("Event.Less() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEvent_Equal(t *testing.T) {
	type fields struct {
		Summary           string
		User              string
		Start             string
		End               string
		Location          string
		CreationTimestamp int64
		CreationDate      string
		UUID              string
	}
	type args struct {
		to *Event
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Event{
				Summary:      tt.fields.Summary,
				User:         tt.fields.User,
				Start:        tt.fields.Start,
				End:          tt.fields.End,
				Location:     tt.fields.Location,
				Timestamp:    tt.fields.CreationTimestamp,
				CreationDate: tt.fields.CreationDate,
				UUID:         tt.fields.UUID,
			}
			if got := e.Equal(tt.args.to); got != tt.want {
				t.Errorf("Event.Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_improbableButMaybeTheyHaveTheSameCreationDate(t *testing.T) {
	type args struct {
		eventI *Event
		eventJ *Event
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"empty", args{&Event{}, &Event{}}, false},
		{"equal", args{&Event{UUID:"97a65b83-ea3f-4fcf-86e8-9defcaafc979"}, &Event{UUID:"97a65b83-ea3f-4fcf-86e8-9defcaafc979"}}, false},
		{"digit v letter", args{&Event{UUID: "16ecb152-e971-470f-9525-a75f380e787b"}, &Event{UUID: "a7e677dd-7f87-49c8-b23b-a4f660e514db"}}, true},
		{"digit v letter reversed", args{ &Event{UUID: "a7e677dd-7f87-49c8-b23b-a4f660e514db"}, &Event{UUID: "16ecb152-e971-470f-9525-a75f380e787b"}}, false},
		{"letter v letter", args{&Event{UUID: "c306bfe1-744f-4868-b959-670d2b821499"}, &Event{UUID: "b6cf408d-e751-4ae7-92cd-b0c522244b77"}}, false},
		{"letter v letter reversed", args{ &Event{UUID: "b6cf408d-e751-4ae7-92cd-b0c522244b77"}, &Event{UUID: "c306bfe1-744f-4868-b959-670d2b821499"}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := improbableButMaybeTheyHaveTheSameCreationDate(tt.args.eventI, tt.args.eventJ); got != tt.want {
				t.Errorf("improbableButMaybeTheyHaveTheSameCreationDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

// This test should compare objects.Event and handlers.CreateRequest
func TestIfEventsMatchRequest(t *testing.T) {
	var input = event.CreateRequest{}
	var _ = Event(input)
}

// This test should check conversion between calendar.Event and objects.Event
func TestIfEventsMatchCalendar(t *testing.T) {
	// TODO
}
