package events

import (
	"github.com/llmcontext/gomcp/protocol/mcp"
)

type EventsProcessor interface {
	EventMcpStared()
	EventMcpInitializeResponse(initializeResponse *mcp.JsonRpcResponseInitializeResult)
	EventMcpToolsListResponse(toolsListResponse *mcp.JsonRpcResponseToolsListResult)
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

func (e *Events) EventMcpToolsListResponse(toolsListResponse *mcp.JsonRpcResponseToolsListResult) {
	e.processor.EventMcpToolsListResponse(toolsListResponse)
}

func (e *Events) EventMcpInitializeResponse(initializeResponse *mcp.JsonRpcResponseInitializeResult) {
	e.processor.EventMcpInitializeResponse(initializeResponse)
}
