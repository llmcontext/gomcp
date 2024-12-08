package prompts

import (
	"encoding/json"

	"github.com/llmcontext/gomcp/types"
)

/*
export const PromptMessageSchema = z
  .object({
    role: z.enum(["user", "assistant"]),
    content: z.union([
      TextContentSchema,
      ImageContentSchema,
      EmbeddedResourceSchema,
    ]),
  })
  .passthrough();

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
    // The base64-encoded image data.
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

type PromptGetResultImpl struct {
	Description string        `json:"description"`
	Messages    []interface{} `json:"messages"`
}

func NewPromptGetResult(description string) types.PromptGetResult {
	return &PromptGetResultImpl{
		Description: description,
		Messages:    []interface{}{},
	}
}

func (r *PromptGetResultImpl) AddTextContent(role types.Role, content string) {
	r.Messages = append(r.Messages, map[string]interface{}{
		"role": role,
		"type": "text",
		"text": content,
	})
}

func (r *PromptGetResultImpl) AddJSONTextContent(role types.Role, content interface{}) {
	// let's marshal the content
	contentBytes, err := json.Marshal(content)
	if err != nil {
		panic(err)
	}

	r.Messages = append(r.Messages, map[string]interface{}{
		"role": role,
		"content": map[string]interface{}{
			"type": "text",
			"text": string(contentBytes),
		},
	})
}

func (r *PromptGetResultImpl) AddImageContent(role types.Role, mimeType string, base64Data string) {
	r.Messages = append(r.Messages, map[string]interface{}{
		"role": role,
		"content": map[string]interface{}{
			"type": "resource",
			"resource": map[string]interface{}{
				"mimeType": mimeType,
				"data":     base64Data,
			},
		},
	})
}

func (r *PromptGetResultImpl) AddEmbeddedResourceTextContent(role types.Role, uri string, mimeType string, text string) {
	r.Messages = append(r.Messages, map[string]interface{}{
		"role": role,
		"content": map[string]interface{}{
			"type": "resource",
			"resource": map[string]interface{}{
				"uri":      uri,
				"mimeType": mimeType,
				"text":     text,
			},
		},
	})
}

func (r *PromptGetResultImpl) AddEmbeddedResourceBlobContent(role types.Role, uri string, mimeType string, base64Data string) {
	r.Messages = append(r.Messages, map[string]interface{}{
		"role": role,
		"content": map[string]interface{}{
			"type": "resource",
			"resource": map[string]interface{}{
				"uri":      uri,
				"mimeType": mimeType,
				"blob":     base64Data,
			},
		},
	})
}
