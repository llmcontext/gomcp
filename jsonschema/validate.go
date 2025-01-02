package jsonschema

import (
	"encoding/json"
	"fmt"

	"github.com/invopop/jsonschema"
	"github.com/xeipuuv/gojsonschema"
)

func ValidateJsonSchemaWithBytes(schema *jsonschema.Schema, data []byte) error {
	schemaBytes, _ := json.MarshalIndent(schema, "", "    ")

	schemaLoader := gojsonschema.NewStringLoader(string(schemaBytes))

	documentLoader := gojsonschema.NewStringLoader(string(data))

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return fmt.Errorf("schema validation error: %v", err)
	}

	if !result.Valid() {
		var errorMessages []string
		for _, desc := range result.Errors() {
			errorMessages = append(errorMessages, desc.String())
		}
		return fmt.Errorf("schema validation failed: %v", errorMessages)
	}
	return nil
}

func ValidateJsonSchemaWithObject(schema *jsonschema.Schema, data interface{}) error {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return ValidateJsonSchemaWithBytes(schema, jsonBytes)
}
