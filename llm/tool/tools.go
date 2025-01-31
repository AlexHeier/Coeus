package tool

import (
	"fmt"
	"reflect"
	"strings"
)

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
	newTool.Name = strings.ToLower(name)

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

	if len(args) != f.Type().NumIn() { // check if the number of arguments is correct
		return nil, fmt.Errorf("wrong number of arguments")
	}

	in := make([]reflect.Value, len(args))
	for i, arg := range args {
		expctedTupe := f.Type().In(i)
		argValue := reflect.ValueOf(arg)

		if !argValue.Type().AssignableTo(expctedTupe) { // check if the argument type is correct
			return nil, fmt.Errorf("wrong argument type expexted %s got %s", expctedTupe, argValue.Type())
		}

		in[i] = argValue // set the argument
	}

	result := f.Call(in) // call the function

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
	name = strings.ToLower(name)
	for _, tool := range Tools {
		if tool.Name == name {
			return tool, nil
		}
	}
	return ToolStruct{}, fmt.Errorf("tool not found")
}
