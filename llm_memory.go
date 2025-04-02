package coeus

import (
	"fmt"
	"strings"
	"time"
)

/* Default memory function is MemoryAllMessage.*/
var memory func(c *Conversation) ([]HistoryStruct, error) = MemoryAllMessage
var memArgs []interface{}

/*
MemoryVersion changes the function used for memory management. Default is All messages.

@param newFunc: The new function to use for memory management.
*/
func MemoryVersion(newFunc ...interface{}) {
	if len(newFunc) > 0 {
		if fn, ok := newFunc[0].(func(*Conversation) ([]HistoryStruct, error)); ok {
			memory = fn
			memArgs = newFunc[1:] // Save additional arguments
		} else {
			memory = MemoryAllMessage // Default if type doesn't match
		}
	} else {
		memory = MemoryAllMessage // Default
	}
}

/*
MemoryAllMessage is a function that will use all messages as memory.

@return Array of history objects to use as memory.
*/
func MemoryAllMessage(c *Conversation) ([]HistoryStruct, error) {

	mem := []HistoryStruct{{Role: "system", Content: sp}}

	return append(mem, c.History...), nil
}

/*
MemoryLastMessage is a function that will use the last int x messages as memory.

@extra param: The number of last user messages to use as memory.

@return Array of the last X amount of messages from user and everything between.
*/
func MemoryLastMessage(c *Conversation) ([]HistoryStruct, error) {
	count, ok := memArgs[0].(int)
	if !ok {
		return nil, fmt.Errorf("memory: second argument needs to be an integer")
	}

	if count < 0 {
		count = -count
	}

	count++ // Add exstra due to user message being appended before memory call.

	i := len(c.History) - 1
	found := 0
	for ; i > 0; i-- {
		if strings.ToLower(c.History[i].Role) == "user" {
			found++
			if found >= count {
				break
			}
		}
	}

	mem := []HistoryStruct{{Role: "system", Content: sp}}

	// Always returns the system message first then the other interactions
	return append(mem, c.History[i:]...), nil
}

/*
MemoryTime is a function that will use the messages within the last int x minutes as memory.

@extra param: The number of last minutes to use as memory.

@return Array of the last X amount of messages within the last Y minutes
*/
func MemoryTime(c *Conversation) ([]HistoryStruct, error) {

	age, ok := memArgs[0].(int)
	if !ok {
		return nil, fmt.Errorf("memory: second argument needs to be an integer")
	}

	if age < 0 {
		age = -age
	}

	mem := []HistoryStruct{{Role: "system", Content: sp}}

	for i := range len(c.History) {
		if time.Since(c.History[i].TimeStamp).Minutes() < float64(age) {
			mem = append(mem, c.History[i])
		}
	}

	return mem, nil
}
