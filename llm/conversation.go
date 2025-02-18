package llm

import (
	"Coeus/llm/tool"
	"Coeus/provider"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"
)

// Struct for containing all the conversations
var ConvAll struct {
	Mutex         sync.Mutex
	Conversations []*Conversation
}

// Struct for a single conversation.
type Conversation struct {
	Mutex      sync.Mutex
	MainPrompt string
	ToolsResp  []interface{}
	History    []string
	Summary    string
	UserPrompt string
	LastActive time.Time
}

type SystemPrompt struct {
	Context struct {
		SystemPrompt string `json:"systemPrompt"`
		Tools        []struct {
			ToolName        string `json:"toolName"`
			ToolDescription string `json:"toolDescription"`
		} `json:"tools"`
		AboutTools  string        `json:"aboutTools"`
		ToolReturns []interface{} `json:"toolReturns"`
		History     []string      `json:"history"`
	} `json:"context"`
}

// Appends a prompt and section to the history within the conversation
func (c *Conversation) AppendHistory(user, llm string) {
	newHistory := "[USER:] " + user + " [LLM:] " + llm
	c.History = append(c.History, newHistory)
}

func (c *Conversation) Prompt(userPrompt string) (provider.ResponseStruct, error) {
	var toolDesc string

	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	c.LastActive = time.Now()
	c.UserPrompt = userPrompt

	if len(tool.Tools) > 0 {
		for _, t := range tool.Tools {
			toolDesc += t.Name + ": " + t.Desc + "\n"
		}
	}

	response, err := provider.Send(provider.RequestStruct{
		Userprompt:   c.UserPrompt,
		Systemprompt: c.BuildSystemPrompt(),
	})
	if err != nil {
		return response, err
	}

	fmt.Println(provider.RequestStruct{
		Userprompt:   c.UserPrompt,
		Systemprompt: c.BuildSystemPrompt(),
	})

	splitString := strings.Split(response.Response, " ")

	// Check for if the response contains a summary and extract it
	for i, w := range splitString {
		if strings.Contains(w, "SUMMARY") {
			c.Summary = strings.Join(splitString[i+1:], " ")
			response.Response = strings.Join(splitString[:i], " ")
			break
		}
	}

	res := response.Response

	for {
		command, _, endIndex := ParseTool(res)
		if command == nil {
			break
		}

		fmt.Println(command)

		args := command[1:]

		t, err := tool.Find(command[0])
		if err != nil {
			fmt.Println(err.Error())
			break
		}

		ft := reflect.ValueOf(t.Function)
		argCount := ft.Type().NumIn()

		callArgs := make([]interface{}, argCount)
		for i := 0; i < argCount; i++ {
			callArgs[i] = args[i]
		}

		tr, err := t.Run(callArgs[0:]...)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		// Gets the last index in function result array which SHOULD be of type error
		err, ok := tr[len(tr)-1].(error)
		if ok {
			fmt.Println(err.Error())
			res = res[endIndex:]
			continue
		}

		c.ToolsResp = append(c.ToolsResp, fmt.Sprintf("%v %v = %v", t.Name, args, tr[:len(tr)-1]))

		response, err = provider.Send(provider.RequestStruct{
			Userprompt:   c.UserPrompt,
			Systemprompt: c.BuildSystemPrompt(),
		})
		if err != nil {
			return response, err
		}
		break
	}
	c.AppendHistory(userPrompt, response.Response)

	return response, nil
}

/*
Parses the first found command within a string by checking if it contains the correct name and format
@param []string: Contains the command name and its args
@param int: Beginning index of where command begins inside string
@param int: End index of where command ends inside string
*/
func ParseTool(s string) ([]string, int, int) {
	s = strings.ReplaceAll(s, "\n", " ")

	for _, t := range tool.Tools {

		beginIndex := strings.Index(s, t.Name)
		var endIndex int
		if beginIndex >= 0 {
			s = s[beginIndex:]

			ft := reflect.ValueOf(t.Function)
			argCount := ft.Type().NumIn()

			var v int
			for end, r := range s {
				if r == ' ' || end == len(s) {
					v++
				}

				if v > argCount {
					s = s[:end]
					endIndex = end
					break
				}
			}
			return strings.Split(s, " "), beginIndex, (endIndex + beginIndex - 1)
		}
	}
	return nil, -1, -1
}

func (c *Conversation) DumpConversation() string {
	var temp string
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	for _, h := range c.History {
		temp += h
	}

	return temp
}

func BeginConversation() *Conversation {
	ConvAll.Mutex.Lock()
	defer ConvAll.Mutex.Unlock()
	newCon := Conversation{
		MainPrompt: Persona,
	}

	ConvAll.Conversations = append(ConvAll.Conversations, &newCon)
	return ConvAll.Conversations[len(ConvAll.Conversations)-1]
}

// Finds and deletes the given conversation
func DeleteConversation(con *Conversation) error {
	ConvAll.Mutex.Lock()
	defer ConvAll.Mutex.Unlock()
	var found bool
	var temp []*Conversation
	for i := range ConvAll.Conversations {
		if !(con == ConvAll.Conversations[i]) {
			temp = append(temp, con)
		} else {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("conversation not found")
	}

	ConvAll.Conversations = temp
	return nil
}

func (c *Conversation) BuildSystemPrompt() string {

	sysP := SystemPrompt{}

	sysP.Context.History = append(sysP.Context.History, c.History...)
	sysP.Context.SystemPrompt = c.MainPrompt
	for _, t := range tool.Tools {
		sysP.Context.Tools = append(sysP.Context.Tools, struct {
			ToolName        string "json:\"toolName\""
			ToolDescription string "json:\"toolDescription\""
		}{ToolName: t.Name, ToolDescription: t.Desc})
	}

	sysP.Context.AboutTools = "Always use tools before using information from your history. To call a tool simply respond with the tool name in all capital letters and its arguments after without any brackets. Example THISISATOOL ARG1 ARG2 ARG... When asked about tool results only answer with the result without brackets and nothing else."
	sysP.Context.ToolReturns = c.ToolsResp

	ret, err := json.MarshalIndent(sysP, "", " ")
	if err != nil {
		fmt.Println(err.Error())
	}

	return string(ret)
}
