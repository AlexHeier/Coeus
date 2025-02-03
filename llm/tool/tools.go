package tool

import (
	"fmt"
	"reflect"
	"strings"
)

// ToolDefintion defines how the LLM should access the tools
var ToolDefintion = `To get access to the TOOLS resources. Respond with the tool name capitalized and the parameters needed. For example, "MULTIPLY 40 60 " `

// ToolStruct is the struct of a tool
type ToolStruct struct {
	Name     string
	Desc     string
	Function interface{}
}

// Tools is a list of all the tools
var Tools []ToolStruct

/*
New creates a new tool and gives access to llm to use it.

@param name: the name of the tool
@param desc: the description of the tool
@param function: the function the llm will call

Sensetive data should not be passed to the function as the LLM will have access to it. Add it as an variable for the thread to use.
*/
func New(name, desc string, function interface{}) {

	var newTool ToolStruct
	newTool.Desc = desc
	newTool.Function = function
	newTool.Name = strings.ToUpper(name)

	Tools = append(Tools, newTool)
}

/*
Run runs the function of the tool its called upon.

@param args: the arguments of the function
@return: the result of the function
@return: an error if the function fails
*/
func (t *ToolStruct) Run(args ...interface{}) ([]interface{}, error) {
	f := reflect.ValueOf(t.Function)
	if f.Kind() != reflect.Func { // check if the function is a function
		return nil, fmt.Errorf("function is not a function")
	}

	// Check if args are an array (or slice) and unpack accordingly
	var finalArgs []reflect.Value

	for _, arg := range args {
		argValue := reflect.ValueOf(arg)

		// If the argument is a slice, we need to process each element in the slice
		if argValue.Kind() == reflect.Slice {
			for j := 0; j < argValue.Len(); j++ {
				elem := argValue.Index(j)
				// Check the type of the element and ensure it is compatible with the function signature
				expectedType := f.Type().In(len(finalArgs))
				if !elem.Type().ConvertibleTo(expectedType) {
					return nil, fmt.Errorf("wrong argument type: expected %s, got %s", expectedType, elem.Type())
				}
				finalArgs = append(finalArgs, elem.Convert(expectedType)) // append converted element
			}
		} else {
			// If it's a regular argument, just check and convert
			expectedType := f.Type().In(len(finalArgs))
			if !argValue.Type().ConvertibleTo(expectedType) {
				return nil, fmt.Errorf("wrong argument type: expected %s, got %s", expectedType, argValue.Type())
			}
			finalArgs = append(finalArgs, argValue.Convert(expectedType))
		}
	}

	// Call the function with the final arguments
	result := f.Call(finalArgs)

	// Convert the result back to []interface{}
	out := make([]interface{}, len(result))
	for i, r := range result {
		out[i] = r.Interface()
	}

	return out, nil // return the result
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
