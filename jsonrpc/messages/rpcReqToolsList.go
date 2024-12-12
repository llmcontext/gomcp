package messages

const (
	RpcRequestMethodToolsList = "tools/list"
)

type JsonRpcRequestToolsListParams struct {
	Cursor *string `json:"cursor,omitempty"`
}
