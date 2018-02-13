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
)

func (em *EventModified) Flag(mod EventModification) (*EventModified) {
	em.Modifications[mod] = eventExists
	return em
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

