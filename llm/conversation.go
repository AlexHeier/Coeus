package llm

import (
	"Coeus/llm/tool"
	"Coeus/provider"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"
)

var ConDB ConversationDB

// Struct for containing all the conversations
type ConversationDB struct {
	M             sync.Mutex
	Conversations []*Conversation
}

// Struct for containing the individual conversations with the LLMs
type Conversation struct {
	M          sync.Mutex
	MainPrompt string
	ToolsResp  interface{}
	History    []string
	Summary    string
	UserPrompt string
	LastActive time.Time
}

// Appends a prompt and section to the history within the conversation
func (c *Conversation) AppendHistory(user, llm string) {
	newHistory := "[USER:] " + user + "[LLM:] " + llm
	c.History = append(c.History, newHistory)
}

func (c *Conversation) Prompt(userPrompt string) (provider.ResponseStruct, error) {
	var toolDesc string
	var response provider.ResponseStruct
	var err error
	var toolResponse []interface{}
	var toolUsed bool = false

	c.M.Lock()
	defer c.M.Unlock()
	c.LastActive = time.Now()

	if len(tool.Tools) > 0 {
		for _, t := range tool.Tools {
			toolDesc += t.Name + ": " + t.Desc + "\n"
		}
	}

	for {
		toolUsed = false

		// Memory(append([]interface{}{c}, MemArgs...)...) sends the conversation and the arguments to the memory function if the user defined some.
		prefix := c.MainPrompt + "[BEGIN TOOLS] Always use a tool if its suitable " + tool.ToolDefintion + toolDesc + "Do not send the tool name if you do not need it. Do not reuse the tool name [END TOOLS] \n [BEGIN TOOL RESPONSE] Answers from used tools, always use these resoults " + fmt.Sprintf("%v", toolResponse) + " [END TOOL RESPONSE]\n" + Memory(append([]interface{}{c}, MemArgs...)...) + "\n\n"

		response, err = provider.Send(prefix + userPrompt)
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

	fmt.Println(userPrompt)
	fmt.Println(response.Response)
	c.AppendHistory(userPrompt, response.Response)
	return response, err
}

func (c *Conversation) DumpConversation() string {
	var temp string
	c.M.Lock()
	defer c.M.Unlock()

	for _, h := range c.History {
		temp += h
	}

	return temp
}

func BeginConversation() *Conversation {
	newCon := Conversation{
		MainPrompt: Persona + "Answer in the language the user is using.\n",
	}

	ConDB.Conversations = append(ConDB.Conversations, &newCon)
	return ConDB.Conversations[len(ConDB.Conversations)-1]
}
