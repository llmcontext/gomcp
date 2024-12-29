package types

type Role string

const (
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

type PromptGetResult interface {
	SetDescription(description string)
	AddTextContent(role Role, content string)
	AddJSONTextContent(role Role, content interface{})
	AddImageContent(role Role, mimeType string, base64Data string)
	AddEmbeddedResourceTextContent(role Role, uri string, mimeType string, text string)
	AddEmbeddedResourceBlobContent(role Role, uri string, mimeType string, base64Data string)
}
