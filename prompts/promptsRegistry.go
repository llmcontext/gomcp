package prompts

type PromptsRegistry struct {
	prompts []PromptDefinition
}

func NewEmptyPromptsRegistry() *PromptsRegistry {
	return &PromptsRegistry{prompts: []PromptDefinition{}}
}

func NewPromptsRegistry(promptYamlFilePath string) (*PromptsRegistry, error) {
	prompts, err := loadPrompts(promptYamlFilePath)
	if err != nil {
		return nil, err
	}
	return &PromptsRegistry{prompts: prompts.Prompts}, nil
}

func (r *PromptsRegistry) GetListOfPrompts() []PromptDefinition {
	return r.prompts
}

func (r *PromptsRegistry) GetPrompt(name string) *PromptDefinition {
	for _, prompt := range r.prompts {
		if prompt.Name == name {
			return &prompt
		}
	}
	return nil
}
