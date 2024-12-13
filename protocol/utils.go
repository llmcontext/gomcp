package protocol

import "fmt"

func CheckIsObject(result interface{}, name string) (map[string]interface{}, error) {
	if result == nil {
		return nil, fmt.Errorf("missing property %s", name)
	}

	resultMap, ok := result.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("%s must be an object", name)
	}
	return resultMap, nil
}

func GetStringField(result map[string]interface{}, name string) (string, error) {
	field, ok := result[name].(string)
	if !ok {
		return "", fmt.Errorf("missing property %s", name)
	}
	return field, nil
}

func GetOptionalStringField(result map[string]interface{}, name string) *string {
	field, ok := result[name].(string)
	if !ok {
		return nil
	}
	return &field
}

func GetBoolField(result map[string]interface{}, name string) (bool, error) {
	field, ok := result[name].(bool)
	if !ok {
		return false, fmt.Errorf("missing property %s", name)
	}
	return field, nil
}

func GetOptionalBoolField(result map[string]interface{}, name string) *bool {
	field, ok := result[name].(bool)
	if !ok {
		return nil
	}
	return &field
}

func GetObjectField(result map[string]interface{}, name string) (map[string]interface{}, error) {
	field, ok := result[name].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("missing property %s", name)
	}
	return field, nil
}

func GetOptionalObjectField(result map[string]interface{}, name string) map[string]interface{} {
	field, ok := result[name].(map[string]interface{})
	if !ok {
		return nil
	}
	return field
}

func GetArrayField(result map[string]interface{}, name string) ([]interface{}, error) {
	field, ok := result[name].([]interface{})
	if !ok {
		return nil, fmt.Errorf("missing property %s", name)
	}
	return field, nil
}

func GetOptionalArrayField(result map[string]interface{}, name string) []interface{} {
	field, ok := result[name].([]interface{})
	if !ok {
		return nil
	}
	return field
}
