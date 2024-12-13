package mcp

import (
	"github.com/llmcontext/gomcp/jsonrpc"
)

type JsonRpcResponsePromptsListResult struct {
	Prompts    []PromptDescription `json:"prompts"`
	NextCursor *string             `json:"nextCursor,omitempty"`
}

type PromptDescription struct {
	Name        string                      `json:"name"`
	Description string                      `json:"description"`
	Arguments   []PromptArgumentDescription `json:"arguments"`
}

type PromptArgumentDescription struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
}

func ParseJsonRpcResponsePromptsList(response *jsonrpc.JsonRpcResponse) (*JsonRpcResponsePromptsListResult, error) {
	resp := JsonRpcResponsePromptsListResult{}

	// parse params
	result, err := checkIsObject(response.Result, "result")
	if err != nil {
		return nil, err
	}

	// read tools
	prompts, err := getArrayField(result, "prompts")
	if err != nil {
		return nil, err
	}

	for _, item := range prompts {
		prompt, err := checkIsObject(item, "prompt")
		if err != nil {
			return nil, err
		}
		name, err := getStringField(prompt, "name")
		if err != nil {
			return nil, err
		}

		description, err := getStringField(prompt, "description")
		if err != nil {
			return nil, err
		}

		arguments, err := getArrayField(prompt, "arguments")
		if err != nil {
			return nil, err
		}

		promptArguments := make([]PromptArgumentDescription, 0)
		for _, argument := range arguments {
			argument, err := checkIsObject(argument, "argument")
			if err != nil {
				return nil, err
			}
			name, err := getStringField(argument, "name")
			if err != nil {
				return nil, err
			}
			description, err := getStringField(argument, "description")
			if err != nil {
				return nil, err
			}
			required, err := getBoolField(argument, "required")
			if err != nil {
				return nil, err
			}
			promptArguments = append(promptArguments, PromptArgumentDescription{
				Name:        name,
				Description: description,
				Required:    required,
			})
		}

		resp.Prompts = append(resp.Prompts, PromptDescription{
			Name:        name,
			Description: description,
			Arguments:   promptArguments,
		})
	}

	// read next cursor
	nextCursor := getOptionalStringField(result, "nextCursor")
	resp.NextCursor = nextCursor

	return &resp, nil
}
