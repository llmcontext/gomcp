package mcp

import (
	"fmt"

	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/protocol"
)

const (
	RpcNotificationMethodResourcesUpdated = "notifications/resources/updated"
)

type JsonRpcNotificationResourcesUpdatedParams struct {
	Uri string `json:"uri"`
}

func ParseJsonRpcNotificationResourcesUpdatedParams(params *jsonrpc.JsonRpcParams) (*JsonRpcNotificationResourcesUpdatedParams, error) {
	if params == nil {
		return nil, fmt.Errorf("invalid call parameters, not an object")
	}
	if !params.IsNamed() {
		return nil, fmt.Errorf("params must be an object")
	}
	namedParams := params.NamedParams

	uri, err := protocol.GetStringField(namedParams, "uri")
	if err != nil {
		return nil, fmt.Errorf("uri is required")
	}

	resourceUpdated := &JsonRpcNotificationResourcesUpdatedParams{
		Uri: uri,
	}

	return resourceUpdated, nil
}
