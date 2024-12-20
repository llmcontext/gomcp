package events

type EventsProcessor interface {
	EventNewProxyTools()
}

type Events struct {
	processor EventsProcessor
}

func NewEvents(processor EventsProcessor) *Events {
	return &Events{
		processor: processor,
	}
}

func (e *Events) EventNewProxyTools() {
	e.processor.EventNewProxyTools()
}
