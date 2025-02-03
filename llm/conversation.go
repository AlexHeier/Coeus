package llm

import (
	"Coeus/llm/tool"
	"Coeus/provider"
	"fmt"
	"reflect"
	"strings"
)

var conversations []Conversation

// Struct for containing the individual conversations with the LLMs
type Conversation struct {
	MainPrompt string
	ToolsResp  interface{}
	History    []string
	Summary    string
	UserPrompt string
}

// Appends a prompt and section to the history within the conversation
func (c *Conversation) AppendHistory(user, llm string) {
	newHistory := "\n[USER:]" + user + "\n[LLM:]" + llm
	c.History = append(c.History, newHistory)
}

func (c *Conversation) Prompt(s string) (provider.ResponseStruct, error) {
	var toolDesc string
	var resString string
	var response provider.ResponseStruct
	var err error
	var toolResponse []interface{}

	if len(tool.Tools) > 0 {
		for _, t := range tool.Tools {
			toolDesc += t.Name + ": " + t.Desc + "\n"
		}
	}

	for {
		var toolUsed bool = false
		// Memory(append([]interface{}{c}, MemArgs...)...) sends the conversation and the arguments to the memory function if the user defined some.
		prefix := c.MainPrompt + "[BEGIN TOOLS] " + tool.ToolDefintion + toolDesc + "[END TOOLS]\n[BEGIN INFORMATION]" + fmt.Sprintf("%v", toolResponse) + "[END INFORMATION]\n[BEGIN HISTORY]" + Memory(append([]interface{}{c}, MemArgs...)...) + "[END HISTORY]\n"
		response, err = provider.Send(prefix + s)
		if err != nil {
			return response, err
		}

		resString = response.Response

		// Check for if the response contains a summary and extract it
		if strings.Contains(resString, "[BEGIN SUMMARY]") {
			beginIndex := strings.Index(resString, "[BEGIN SUMMARY]")
			endIndex := strings.Index(resString, "[END SUMMARY]")
			if endIndex > beginIndex {
				c.Summary = strings.TrimSpace(resString[beginIndex:endIndex])
			}
		}

		for _, t := range tool.Tools {
			if strings.Contains(resString, t.Name) {
				toolUsed = true
				ft := reflect.ValueOf(t.Function)
				argCount := ft.Type().NumIn()

				beginIndex := strings.Index(resString, t.Name)
				args := strings.Fields(resString[beginIndex:])
				args = args[:argCount]

				callArgs := make([]reflect.Value, argCount)
				for i := 0; i < argCount; i++ {
					callArgs[i] = reflect.ValueOf(args[i])
				}
				tr, err := t.Run(callArgs)
				if err != nil {
					return response, err
				}
				toolResponse = append(toolResponse, tr)
			}
		}
		if !toolUsed {
			break
		}
	}

	c.AppendHistory(s, resString)
	fmt.Println(response)
	return response, err
}

func (c *Conversation) DumpConversation() string {
	dump := c.MainPrompt

	for _, h := range c.History {
		dump += h
	}
	return dump
}

func BeginConversation() *Conversation {
	newCon := Conversation{
		MainPrompt: Persona + "Answer in the language the user is using.\n",
	}

	conversations = append(conversations, newCon)
	return &conversations[len(conversations)-1]
}
