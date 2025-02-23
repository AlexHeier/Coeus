package llm

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/AlexHeier/Coeus/provider"
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
	History    []provider.HistoryStruct
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

func (c *Conversation) AppendHistory(role, content string) {
	c.History = append(c.History, provider.HistoryStruct{Role: role, Content: content})
}

func (c *Conversation) Prompt(userPrompt string) (provider.ResponseStruct, error) {

	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	c.LastActive = time.Now()
	c.UserPrompt = userPrompt

	response, err := provider.Send(provider.RequestStruct{
		Userprompt:   c.UserPrompt,
		Systemprompt: c.systemPrompt(),
		History:      &c.History,
	})

	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("wtf")
		return response, err
	}

	c.AppendHistory("user", c.UserPrompt)
	c.AppendHistory("assistant", response.Response)

	splitString := strings.Split(response.Response, " ")

	//Check for if the response contains a summary and extract it

	for i, w := range splitString {
		if strings.Contains(w, "SUMMARY") {
			c.Summary = strings.Join(splitString[i+1:], " ")
			response.Response = strings.Join(splitString[:i], " ")
			break

		}
	}
	return response, err
}

func (c *Conversation) DumpConversation() string {
	var temp string
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	for _, h := range c.History {
		temp += h.Role + ": " + h.Content + "\n"
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

func (c *Conversation) systemPrompt() string {

	sysP := make(map[string]interface{})

	sysP["systemprompt"] = c.MainPrompt
	sysP["conversationHistory"] = Memory(append([]interface{}{c}, MemArgs...)...)

	ret, err := json.Marshal(sysP)
	if err != nil {
		fmt.Println(err.Error())
	}

	return string(ret)
}
