package objects

func NewModified(event *Event) (*EventModified) {
	return &EventModified{
		Event:         event,
		Modifications: make(map[EventModification]struct{}),
	}
}

type EventModified struct {
	Event         *Event
	Modifications map[EventModification]struct{}
}
var eventExists = struct{}{}

type EventModification uint

const (
	Deleted          EventModification = iota
	Added
	ModifiedTime
	ModifiedLocation
	UserAccepted
	UserRejected

	SomethingWonkyHappened
)

func (em *EventModified) Flag(mod EventModification) (*EventModified) {
	em.Modifications[mod] = eventExists
	return em
}

func (em *EventModified) ToListOfWords() ([]string) {
	words := make([]string, 0)
	for k := range em.Modifications {
		var word string

		switch k {
		case Deleted:
			word = "deleted"
		case Added:
			word = "added"
		case ModifiedTime:
			word = "modified time"
		case ModifiedLocation:
			word = "modified location"
		case UserAccepted:
			word = "user accepted"
		case UserRejected:
			word = "user rejected"
		case SomethingWonkyHappened:
			word = "something wonky happened"
		}

		words = append(words, word)
	}
	return words
}

// unnecessary, but leaving this just in case

//type SortableModifiedEvents []*EventModified
//
//func (sme SortableModifiedEvents) Len() int {
//	return len(sme)
//}
//
//func (sme SortableModifiedEvents) Less(i, j int) bool {
//	sme[i].Event.Less(sme[j].Event)
//}
//
//func (sme SortableModifiedEvents) Swap(i, j int) {
//	sme[i], sme[j] = sme[j], sme[i]
//}

