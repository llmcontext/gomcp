package prompts

// PromptList represents the root of a yaml file containing a list of prompts
type PromptList struct {
	Prompts []*PromptDefinition `json:"prompts" yaml:"prompts"`
}

type PromptDefinition struct {
	Name        string                     `json:"name" yaml:"name"`
	Description string                     `json:"description" yaml:"description"`
	Arguments   []PromptArgumentDefinition `json:"arguments,omitempty" yaml:"arguments,omitempty"`
	Prompt      string                     `json:"prompt" yaml:"prompt"`
}

type PromptArgumentDefinition struct {
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description" yaml:"description"`
	Required    bool   `json:"required" yaml:"required"`
}

type DuplicatedPrompt struct {
	PromptName string
	FilePath   string
}
