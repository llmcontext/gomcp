package types

type ToolCallResult interface {
	AddTextContent(content string)
	AddJSONTextContent(content interface{})
	AddImageContent(base64Data string, mimeType string)
	AddEmbeddedResourceTextContent(uri string, mimeType string, text string)
	AddEmbeddedResourceBlobContent(uri string, mimeType string, base64Data string)
	SetError(isError bool)
}
