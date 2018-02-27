package objects

import (
	"calendar-synch/helpers"
)

type Event struct {
	Summary      string `json:"summary"`
	User         string `json:"user"`
	Start        string `json:"start"`
	End          string `json:"end"`
	Location     string `json:"location"`
	CreationDate string `json:"creationDate"`
	UUID         string `json:"-"`
}

//IsTheSame checks if two events have the same fields (used in checking if they've been changed).
//This works only if the struct doesn't have any slices.
func (e *Event) IsTheSame(to *Event) (bool) {
	return *e == *to
}

//Less compares creation date
func (e *Event) Less(than *Event) (bool) {
	left := helpers.StringToTime(e.CreationDate)
	right := helpers.StringToTime(than.CreationDate)

	// just in fucking case
	if left.Equal(right) {
		return improbableButMaybeTheyHaveTheSameCreationDate(e, than)
	}

	return left.Before(right)
}

func (e *Event) Equal(to *Event) (bool) {
	left := helpers.StringToTime(e.CreationDate)
	right := helpers.StringToTime(to.CreationDate)
	return left.Equal(right)
}

// FIXME actually this is very much probable as the creation date is precise up to 1 second
//and it's also a pity to throw away such a "beautiful" function xd
func improbableButMaybeTheyHaveTheSameCreationDate(eventI, eventJ *Event) (bool) {
	si := helpers.StringToTime(eventI.Start)
	ei := helpers.StringToTime(eventI.End)
	sj := helpers.StringToTime(eventJ.Start)
	ej := helpers.StringToTime(eventJ.End)

	// si,ei,sj,ej come from SortableEvents.Less()
	// start_of_i'th, end_of_i'th, etc ...

	// ---si----ei-->  this one is smaller
	// ----sj---ei-->

	// ---si--ei---->  this one is smaller
	// ---sj---ei--->

	// ----si-ei---->  added few weeks later: TODO what about this one?
	// ---sj---ei--->

	if si.Equal(sj) {
		return ei.Before(ej)
	}
	return si.Before(sj)
}

type SortableEvents []*Event

func (s SortableEvents) Len() int {
	return len(s)
}

// slowest sort in existence
func (s SortableEvents) Less(i, j int) bool {
	//
	return s[i].Less(s[j])
}

func (s SortableEvents) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// TODO ancestors
