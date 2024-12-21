package mcp

type JsonRpcResponsePromptsGetResult struct {
	Description string          `json:"description"`
	Messages    []PromptMessage `json:"messages"`
}

type PromptMessage struct {
	Role    string `json:"role"` // "user" or "assistant"
	Content string `json:"content"`
}
