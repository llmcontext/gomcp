package jsonschema

import (
	"encoding/json"
	"fmt"

	"github.com/invopop/jsonschema"
)

func ToJsonSchema(rawschema interface{}) (*jsonschema.Schema, error) {
	var schema *jsonschema.Schema
	switch s := rawschema.(type) {
	case *jsonschema.Schema:
		schema = s
	case map[string]interface{}:
		schema = &jsonschema.Schema{}
		// Unmarshal the map into the schema
		if err := mapToStruct(s, schema); err != nil {
			return nil, fmt.Errorf("invalid schema format: %v", err)
		}
	default:
		return nil, fmt.Errorf("inputSchema must be either *jsonschema.Schema or map[string]interface{}")
	}
	return schema, nil
}

func mapToStruct(input map[string]interface{}, output interface{}) error {
	jsonBytes, err := json.Marshal(input)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonBytes, output)
}
