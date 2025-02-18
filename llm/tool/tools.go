package tool

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

type ToolStruct struct {
	Name     string                 `json:"name"`
	Desc     string                 `json:"description"`
	Params   map[string]interface{} `json:"parameters"`
	Function interface{}            `json:"-"` // Function execution
}

// Tools is a list of all the tools
var Tools []ToolStruct

func New(name, desc string, function interface{}) {
	paramSchema := extractFunctionParams(function)

	newTool := ToolStruct{
		Name:     strings.ToUpper(name),
		Desc:     desc,
		Params:   paramSchema,
		Function: function,
	}

	Tools = append(Tools, newTool)
}

/*
Run runs the function of the tool its called upon.

@param args: the arguments of the function
@return: the result of the function
@return: an error if the function fails
*/
func (t *ToolStruct) Run(args ...interface{}) (string, error) {
	f := reflect.ValueOf(t.Function)
	if f.Kind() != reflect.Func {
		return "", fmt.Errorf("function is not a function")
	}

	// Check if args are an array (or slice) and unpack accordingly
	var finalArgs []reflect.Value

	for _, arg := range args {
		argValue := reflect.ValueOf(arg)

		// If the argument is a slice, process each element
		if argValue.Kind() == reflect.Slice {
			for j := 0; j < argValue.Len(); j++ {
				elem := argValue.Index(j)
				expectedType := f.Type().In(len(finalArgs))
				if !elem.Type().ConvertibleTo(expectedType) {
					return "", fmt.Errorf("wrong argument type: expected %s, got %s", expectedType, elem.Type())
				}
				finalArgs = append(finalArgs, elem.Convert(expectedType))
			}
		} else {
			// If it's a regular argument, check and convert
			expectedType := f.Type().In(len(finalArgs))
			if !argValue.Type().ConvertibleTo(expectedType) {
				return "", fmt.Errorf("wrong argument type: expected %s, got %s", expectedType, argValue.Type())
			}
			finalArgs = append(finalArgs, argValue.Convert(expectedType))
		}
	}

	// Call the function with the final arguments
	result := f.Call(finalArgs)

	// Convert the result to a JSON response
	responseData := make([]interface{}, len(result))
	for i, r := range result {
		responseData[i] = r.Interface()
	}

	jsonResponse, err := json.Marshal(responseData)
	if err != nil {
		return "", fmt.Errorf("error converting result to JSON: %v", err)
	}

	return string(jsonResponse), nil
}

/*
Find finds a tool by its name and returns the tool struct.

@param name: the name of the tool
@return: the tool struct
@return: an error if the tool is not found
*/
func Find(name string) (ToolStruct, error) {
	for _, tool := range Tools {
		if tool.Name == name {
			return tool, nil
		}
	}
	return ToolStruct{}, fmt.Errorf("tool not found")
}

/*
Returns the information about each tool to be used by an LLM. For use within a prompt

@return: Nothing if no tools specified
@return: Tool names and descriptions with usage information
*/
func GetToolsDescription() string {

	if len(Tools) <= 0 {
		return ""
	}

	var desc string

	for _, tool := range Tools {
		desc += tool.Name + ": " + tool.Desc + " "
	}

	return desc
}

func extractFunctionParams(fn interface{}) map[string]interface{} {
	fnType := reflect.TypeOf(fn)
	if fnType.Kind() != reflect.Func {
		panic("Function must be of type func")
	}

	params := make(map[string]interface{})
	properties := make(map[string]interface{})
	var required []string

	// Extract parameter information
	for i := 0; i < fnType.NumIn(); i++ {
		paramName := fmt.Sprintf("param%d", i+1)
		paramType := fnType.In(i).Kind().String()

		if paramType == "int" {
			paramType = "number"
		}

		properties[paramName] = map[string]string{
			"type": paramType,
		}
		required = append(required, paramName)
	}

	params["type"] = "object"
	params["properties"] = properties
	if len(required) > 0 {
		params["required"] = required
	}

	return params
}
