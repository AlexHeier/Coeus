package llm

import (
	"Coeus/provider"
	"encoding/json"
	"fmt"
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

// Appends a prompt and section to the history within the conversation
func (c *Conversation) AppendHistory(user, llm string) {
	newHistory := "[USER:] " + user + "[LLM:] " + llm
	c.History = append(c.History, newHistory)
}

func (c *Conversation) Prompt(userPrompt string) (provider.ResponseStruct, error) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	c.LastActive = time.Now()
	c.UserPrompt = userPrompt

	response, err := provider.Send(provider.RequestStruct{
		Userprompt:   c.UserPrompt,
		Systemprompt: c.systemPrompt(),
	})
	if err != nil {
		return response, err
	}
	c.AppendHistory(userPrompt, response.Response)

	splitString := strings.Split(response.Response, " ")

	// Check for if the response contains a summary and extract it
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
