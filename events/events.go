package events

type Type int

const (
	Unknown Type = iota
	Message
)

type Event struct {
	Type Type
	Text string
}

type Fetcher interface {
	Fetch(limit int) ([]Event, error)
}

type Processor interface {
	Process(e Event) error
}
