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

func (c *Conversation) Prompt(UserPrompt string) (provider.ResponseStruct, error) {
	var toolDesc string
	var response provider.ResponseStruct
	var err error
	var toolResponse []interface{}
	var toolUsed bool = false

	if len(tool.Tools) > 0 {
		for _, t := range tool.Tools {
			toolDesc += t.Name + ": " + t.Desc + "\n"
		}
	}

	for {
		toolUsed = false
		// Memory(append([]interface{}{c}, MemArgs...)...) sends the conversation and the arguments to the memory function if the user defined some.
		prefix := c.MainPrompt + "Always use a tool if its suitable " + tool.ToolDefintion + toolDesc + "Do not send the tool name if you do not need it. Do not reuse the tool name\nAnswers from used tools(" + fmt.Sprintf("%v", toolResponse) + ")\n" + Memory(append([]interface{}{c}, MemArgs...)...) + "\n\n"

		response, err = provider.Send(prefix + UserPrompt)
		if err != nil {
			return response, err
		}

		splitString := strings.Split(response.Response, " ")

		// Check for if the response contains a summary and extract it
		for i, w := range splitString {
			if strings.Contains(w, "SUMMARY") {
				c.Summary = strings.Join(splitString[i+1:], " ")
				response.Response = strings.Join(splitString[:i], " ")
				break
			}
		}

		for _, t := range tool.Tools {
			if strings.Contains(response.Response, t.Name) {
				var startIndex int
				toolUsed = true
				for i := range splitString {
					if splitString[i] == t.Name {
						startIndex = i + 1
						break
					}
				}

				ft := reflect.ValueOf(t.Function)
				argCount := ft.Type().NumIn()

				args := splitString[startIndex:]
				args = args[:argCount]

				callArgs := make([]interface{}, argCount)
				for i := 0; i < argCount; i++ {
					callArgs[i] = args[i]
				}
				tr, err := t.Run(callArgs[0:]...)
				if err != nil {
					return response, err
				}

				toolResponse = append(toolResponse, fmt.Sprintf("%v %v = %v", t.Name, args, tr))
			}
		}
		if !toolUsed {
			break
		}
	}

	c.AppendHistory(UserPrompt, response.Response)
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
