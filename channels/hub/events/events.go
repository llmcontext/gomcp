package events

type EventsProcessor interface {
	EventMcpStared()
}

type Events struct {
	processor EventsProcessor
}

func NewEvents(processor EventsProcessor) *Events {
	return &Events{
		processor: processor,
	}
}

func (e *Events) EventMcpStarted() {
	e.processor.EventMcpStared()
}
