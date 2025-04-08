package coeus

import (
	"fmt"
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

/*
ConversationTTL is a function that automatically cleans up conversations by checking their Time To Live (TTL).

@param ttl Time in minutes a conversation can be dormant before being deleted.

@return error If ttl is not a positive number.
*/
func ConversationTTL(ttl int) error {
	if ttl < 0 {
		return fmt.Errorf("TTL has to be a positive number")
	}
	go beginCleanup(ttl)
	return nil
}

func beginCleanup(ttl int) {
	for {
		var temp []*Conversation
		ConvAll.Mutex.Lock()
		for _, con := range ConvAll.Conversations {
			if time.Since(con.LastActive) >= time.Duration(ttl)*time.Minute {
				temp = append(temp, con)
			}
		}
		ConvAll.Conversations = temp
		ConvAll.Mutex.Unlock()
		time.Sleep(30 * time.Second)
	}
}
