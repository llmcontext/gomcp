package registry

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/llmcontext/gomcp/pkg/prompts"
	"github.com/llmcontext/gomcp/providers/results"
	"github.com/llmcontext/gomcp/types"
)

type PromptsRegistry struct {
	prompts []*prompts.PromptDefinition
}

func NewPromptsRegistry() *PromptsRegistry {
	return &PromptsRegistry{prompts: []*prompts.PromptDefinition{}}
}

func (r *PromptsRegistry) LoadPromptYamlFile(promptYamlFilePath string) error {
	prompts, err := prompts.LoadPromptYamlFile(promptYamlFilePath)
	if err != nil {
		return err
	}
	// add the prompts to the registry
	r.prompts = append(r.prompts, prompts.Prompts...)
	return nil
}

func (r *PromptsRegistry) GetListOfPrompts() []*prompts.PromptDefinition {
	return r.prompts
}

func (r *PromptsRegistry) findPrompt(name string) *prompts.PromptDefinition {
	for _, prompt := range r.prompts {
		if prompt.Name == name {
			return prompt
		}
	}
	return nil
}

func (r *PromptsRegistry) GetPrompt(promptName string, arguments map[string]string) (interface{}, error) {
	prompt := r.findPrompt(promptName)
	if prompt == nil {
		return nil, fmt.Errorf("prompt %s not found", promptName)
	}

	var templateArgs = make(map[string]string)

	// let's go through all the arguments one by one
	for _, argument := range prompt.Arguments {
		argumentValue, ok := arguments[argument.Name]
		if argument.Required && !ok {
			return nil, fmt.Errorf("missing argument: %s", argument.Name)
		}

		templateArgs[argument.Name] = argumentValue
	}

	tmpl, err := template.New(promptName).Parse(prompt.Prompt)
	if err != nil {
		return nil, fmt.Errorf("invalid prompt template: %s", err)
	}

	var renderedPrompt bytes.Buffer
	err = tmpl.Execute(&renderedPrompt, templateArgs)
	if err != nil {
		return nil, fmt.Errorf("invalid prompt rendering: %s", err)
	}
	promptResult := renderedPrompt.String()

	// let's create the output
	output := results.NewPromptGetResult(prompt.Description)

	output.AddTextContent(types.RoleUser, promptResult)

	return output, nil
}
