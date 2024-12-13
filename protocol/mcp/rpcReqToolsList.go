package mcp

const (
	RpcRequestMethodToolsList = "tools/list"
)

type JsonRpcRequestToolsListParams struct {
	Cursor *string `json:"cursor,omitempty"`
}
