package mcp

import (
	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/protocol"
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
	result, err := protocol.CheckIsObject(response.Result, "result")
	if err != nil {
		return nil, err
	}

	// read tools
	prompts, err := protocol.GetArrayField(result, "prompts")
	if err != nil {
		return nil, err
	}

	for _, item := range prompts {
		prompt, err := protocol.CheckIsObject(item, "prompt")
		if err != nil {
			return nil, err
		}
		name, err := protocol.GetStringField(prompt, "name")
		if err != nil {
			return nil, err
		}

		description, err := protocol.GetStringField(prompt, "description")
		if err != nil {
			return nil, err
		}

		arguments, err := protocol.GetArrayField(prompt, "arguments")
		if err != nil {
			return nil, err
		}

		promptArguments := make([]PromptArgumentDescription, 0)
		for _, argument := range arguments {
			argument, err := protocol.CheckIsObject(argument, "argument")
			if err != nil {
				return nil, err
			}
			name, err := protocol.GetStringField(argument, "name")
			if err != nil {
				return nil, err
			}
			description, err := protocol.GetStringField(argument, "description")
			if err != nil {
				return nil, err
			}
			required, err := protocol.GetBoolField(argument, "required")
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
	nextCursor := protocol.GetOptionalStringField(result, "nextCursor")
	resp.NextCursor = nextCursor

	return &resp, nil
}
