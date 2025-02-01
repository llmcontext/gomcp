package jsonschema

import (
	"fmt"
	"reflect"

	"github.com/invopop/jsonschema"
)

func GetSchemaFromAny(any interface{}) (*jsonschema.Schema, error) {
	// let's generate the schema from the config struct
	jsonSchema := jsonschema.Reflect(any)
	if jsonSchema == nil {
		return nil, fmt.Errorf("failed to generate schema from struct")
	}
	return jsonSchema, nil
}

func GetSchemaFromType(t reflect.Type) (*jsonschema.Schema, string, error) {
	var typeName = t.Elem().Name()
	if typeName == "" {
		return nil, "", fmt.Errorf("type name is empty")
	}

	reflector := jsonschema.Reflector{}
	schema := reflector.ReflectFromType(t)
	if schema == nil {
		return nil, "", fmt.Errorf("error generating schema")
	}

	schemaType := schema.Definitions[typeName]
	if schemaType == nil {
		return nil, "", fmt.Errorf("no schema for definition found")
	}
	return schemaType, typeName, nil
}

func GetFullSchemaFromInterface(t reflect.Type) (*jsonschema.Schema, string, error) {
	var typeName = t.Elem().Name()
	if typeName == "" {
		return nil, "", fmt.Errorf("type name is empty")
	}

	reflector := jsonschema.Reflector{}
	schema := reflector.ReflectFromType(t)
	if schema == nil {
		return nil, "", fmt.Errorf("error generating schema")
	}

	return schema, typeName, nil
}
