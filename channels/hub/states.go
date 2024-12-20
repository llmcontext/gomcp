package hub

import (
	"github.com/llmcontext/gomcp/channels/hub/events"
	"github.com/llmcontext/gomcp/types"
)

type StateManager struct {
	logger types.Logger
}

func NewStateManager(logger types.Logger) *StateManager {
	return &StateManager{
		logger: logger,
	}
}

func (s *StateManager) AsEventsProcessor() events.EventsProcessor {
	return s
}

func (s *StateManager) EventMcpStared() {
	// TODO: implement
}
