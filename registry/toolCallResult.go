package registry

import (
	"encoding/json"

	"github.com/llmcontext/gomcp/types"
)

/*
 response must be:

 https://github.com/modelcontextprotocol/typescript-sdk/blob/main/src/types.ts#L56

 export const CallToolResultSchema = ResultSchema.extend({
  content: z.array(
    z.union([TextContentSchema, ImageContentSchema, EmbeddedResourceSchema]),
  ),
  isError: z.boolean().default(false).optional(),
});

export const TextContentSchema = z
  .object({
    type: z.literal("text"),
	// The text content of the message.
    text: z.string(),
  })
  .passthrough();

export const ImageContentSchema = z
  .object({
    type: z.literal("image"),
    //The base64-encoded image data.
    data: z.string().base64(),
     // The MIME type of the image. Different providers may support different image types.
    mimeType: z.string(),
  })
  .passthrough();

export const EmbeddedResourceSchema = z
  .object({
    type: z.literal("resource"),
    resource: z.union([TextResourceContentsSchema, BlobResourceContentsSchema]),
  })
  .passthrough();

export const ResourceContentsSchema = z
  .object({
    // The URI of this resource.
    uri: z.string(),
    // The MIME type of this resource, if known.
    mimeType: z.optional(z.string()),
  })
  .passthrough();

export const TextResourceContentsSchema = ResourceContentsSchema.extend({
  // The text of the item. This must only be set if the item can actually be represented as text (not binary data).
  text: z.string(),
});

export const BlobResourceContentsSchema = ResourceContentsSchema.extend({
  // A base64-encoded string representing the binary data of the item.
  blob: z.string().base64(),
});
*/

type ToolCallResultImpl struct {
	Content []interface{} `json:"content"`
	IsError *bool         `json:"isError,omitempty"`
}

func NewToolCallResult() types.ToolCallResult {
	return &ToolCallResultImpl{
		Content: []interface{}{},
		IsError: nil,
	}
}

func (r *ToolCallResultImpl) AddTextContent(content string) {
	r.Content = append(r.Content, map[string]interface{}{
		"type": "text",
		"text": content,
	})
}

func (r *ToolCallResultImpl) AddJSONTextContent(content interface{}) {
	// let's marshal the content
	contentBytes, err := json.Marshal(content)
	if err != nil {
		panic(err)
	}
	r.Content = append(r.Content, map[string]interface{}{
		"type": "text",
		"text": string(contentBytes),
	})
}

func (r *ToolCallResultImpl) AddImageContent(mimeType string, base64Data string) {
	r.Content = append(r.Content, map[string]interface{}{
		"type":     "image",
		"data":     base64Data,
		"mimeType": mimeType,
	})
}

func (r *ToolCallResultImpl) AddEmbeddedResourceTextContent(uri string, mimeType string, text string) {
	r.Content = append(r.Content, map[string]interface{}{
		"type": "resource",
		"resource": map[string]interface{}{
			"uri":      uri,
			"mimeType": mimeType,
			"text":     text,
		},
	})
}

func (r *ToolCallResultImpl) AddEmbeddedResourceBlobContent(uri string, mimeType string, base64Data string) {
	r.Content = append(r.Content, map[string]interface{}{
		"type": "resource",
		"resource": map[string]interface{}{
			"uri":      uri,
			"mimeType": mimeType,
			"blob":     base64Data,
		},
	})
}

func (r *ToolCallResultImpl) SetError(isError bool) {
	r.IsError = &isError
}
