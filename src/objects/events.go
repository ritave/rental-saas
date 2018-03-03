package objects

type Event struct {
	Summary  string `json:"summary"`
	User     string `json:"user"`
	Start    string `json:"start"`
	End      string `json:"end"`
	Location string `json:"location"`
	// number of MILLISECONDS (agreeing on Google's terms)
	Timestamp    int64  `json:"timestamp"`
	CreationDate string
	UUID         string `json:"-"`
}

//IsTheSame checks if two events have the same fields (used in checking if they've been changed).
//This works only if the struct doesn't have any slices.
func (e *Event) IsTheSame(to *Event) (bool) {
	return *e == *to
}

//Less compares creation date
func (e *Event) Less(than *Event) (bool) {
	// just in fucking case
	if e.Timestamp == than.Timestamp {
		return improbableButMaybeTheyHaveTheSameCreationDate(e, than)
	}

	return e.Timestamp < than.Timestamp
}

func (e *Event) Equal(to *Event) (bool) {
	return e.Timestamp == to.Timestamp
}

func improbableButMaybeTheyHaveTheSameCreationDate(eventI, eventJ *Event) (bool) {
	return eventI.UUID < eventJ.UUID
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
