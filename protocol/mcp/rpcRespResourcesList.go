package mcp

type JsonRpcResponseResourcesListResult struct {
	Resources []ResourceDescription `json:"resources"`
}
type ResourceDescription struct {
	Uri         string `json:"uri"`
	Name        string `json:"name"`
	Description string `json:"description"`
	MimeType    string `json:"mimeType"`
}
