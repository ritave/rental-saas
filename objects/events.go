package objects

import (
	"calendar-synch/helpers"
)

type Event struct {
	Summary  string
	User     string
	Start    string
	End      string
	Location string
}

//Equal checks if two events are equal (used in checking if they've been changed).
//This works only if the struct doesn't have any slices.
func (e *Event) Equal(to *Event) (bool) {
	return *e == *to
}

type SortableEvents []*Event

func (s SortableEvents) Len() int {
	return len(s)
}

// slowest sort in existence
func (s SortableEvents) Less(i, j int) bool {
	si := helpers.StringToTime(s[i].Start)
	ei := helpers.StringToTime(s[i].End)
	sj := helpers.StringToTime(s[j].Start)
	ej := helpers.StringToTime(s[j].End)

	// ---si----ei-->  this one is smaller
	// ----sj---ei-->

	// ---si--ei---->  this one is smaller
	// ---sj---ei--->

	if si.Equal(sj) {
		return ei.Before(ej)
	}
	return si.Before(sj)
}

func (s SortableEvents) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// TODO ancestors


