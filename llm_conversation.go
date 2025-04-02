package coeus

import (
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
	History    []HistoryStruct
	Summary    string
	UserPrompt string
	LastActive time.Time
}

// Appends a prompt and section to the history within the conversation
func (c *Conversation) appendHistory(role, content string) {
	c.History = append(c.History, HistoryStruct{Role: role, Content: content, TimeStamp: time.Now()})
}

/*
Prompt is a function that sends a prompt to the LLM and returns the response.

@receiver c: The conversation to send the prompt from

@param userPrompt: The prompt to send to the LLM

@return A ResponseStruct and an error if the request fails
*/
func (c *Conversation) Prompt(userPrompt string) (ResponseStruct, error) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	c.LastActive = time.Now()
	c.UserPrompt = userPrompt

	c.appendHistory("user", userPrompt)

	response, err := Send(c)
	if err != nil {
		fmt.Println(err.Error())
		return response, err
	}

	c.appendHistory("assistant", response.Response)

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

/*
DumpConversation is a function that returns the conversation history as a string.

@receiver c: The conversation to dump
*/
func (c *Conversation) DumpConversation() string {
	var temp string
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	for _, h := range c.History {
		temp += h.Role + ": " + h.Content + "\n"
	}
	return temp
}

/*
BeginConversation is a function that creates a new conversation and returns it.

@return A pointer to the new conversation
*/
func BeginConversation() *Conversation {
	ConvAll.Mutex.Lock()
	defer ConvAll.Mutex.Unlock()
	newCon := Conversation{
		MainPrompt: sp,
	}
	ConvAll.Conversations = append(ConvAll.Conversations, &newCon)
	return ConvAll.Conversations[len(ConvAll.Conversations)-1]
}
