package sdk

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// callTool invokes a tool by name with provided arguments
// returns the result, error, error
// the first error is the error from the function when called,
// the second error is if there was an error validating the arguments
func callFunction(fn interface{}, args ...interface{}) (interface{}, error, error) {
	fnValue := reflect.ValueOf(fn)
	fnType := fnValue.Type()
	fnName := fnType.Name()

	// Validate argument count
	if len(args) != fnType.NumIn() {
		return nil, nil, fmt.Errorf("fn %s expects %d arguments, got %d", fnName, fnType.NumIn(), len(args))
	}

	// Convert arguments to reflect.Value slice
	callArgs := make([]reflect.Value, len(args))
	for i, arg := range args {
		expectedType := fnType.In(i)

		// Handle map[string]interface{} to struct conversion
		if mapArg, ok := arg.(map[string]interface{}); ok {
			// Create a new instance of the expected type
			newArg := reflect.New(expectedType).Interface()

			// Convert map to JSON
			jsonBytes, err := json.Marshal(mapArg)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to marshal argument %d: %v", i, err)
			}

			// Unmarshal JSON into the typed struct
			if err := json.Unmarshal(jsonBytes, newArg); err != nil {
				return nil, nil, fmt.Errorf("failed to unmarshal argument %d: %v", i, err)
			}

			callArgs[i] = reflect.ValueOf(newArg).Elem()
		} else {
			// Original type checking for non-map arguments
			if !reflect.TypeOf(arg).AssignableTo(expectedType) {
				fmt.Printf("invalid argument type for parameter %d: %v\n", i, reflect.TypeOf(arg))
				fmt.Printf("expected type: %v\n", expectedType)
				return nil, nil, fmt.Errorf("invalid argument type for parameter %d", i)
			}
			callArgs[i] = reflect.ValueOf(arg)
		}
	}

	// Call the function
	results := fnValue.Call(callArgs)

	var callError error

	if len(results) == 1 {
		// assume error
		if results[0].IsNil() {
			callError = nil
		} else {
			callError = results[0].Interface().(error)
		}
		return nil, callError, nil
	}

	// Handle return values
	if len(results) == 2 { // Assuming (result, error) pattern
		if !results[1].IsNil() {
			callError = results[1].Interface().(error)
		} else {
			callError = nil
		}
		return results[0].Interface(), callError, nil
	}

	return nil, nil, fmt.Errorf("fn %s returned %d values, expected 1 or 2", fnName, len(results))
}
