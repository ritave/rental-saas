package objects

func NewModified(event *Event) (*EventModified) {
	return &EventModified{
		Event:         event,
		Modifications: make([]EventModification, 0),
	}
}

type EventModified struct {
	Event         *Event
	Modifications []EventModification
}

type EventModification uint

const (
	Deleted          EventModification = iota
	Added
	ModifiedTime
	ModifiedLocation
)

func (em *EventModified) Flag(mod EventModification) (*EventModified) {
	em.Modifications = append(em.Modifications, mod)
	return em
}
