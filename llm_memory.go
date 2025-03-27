package coeus

import (
	"fmt"
	"time"
)

/* Default memory function is All. This function will use all messages as memory. */
var memory func(args ...interface{}) ([]HistoryStruct, error) = MemoryAllMessage
var memArgs []interface{}

const memerr string = "memory: argument 1 is not a conversation"

/*
MemoryVersion changes the function used for memory management. Default is All messages.

@param newFunc: The new function to use for memory management.
*/
func MemoryVersion(newFunc ...interface{}) {
	if len(newFunc) > 0 {
		if fn, ok := newFunc[0].(func(args ...interface{}) ([]HistoryStruct, error)); ok {
			memory = fn
			if len(newFunc) > 1 {
				memArgs = newFunc[1:]
			}
		}
	} else {
		memory = MemoryAllMessage // Default
		return
	}
}

/*
MemoryAllMessage is a function that will use all messages as memory.

@return Array of history objects to use as memory.
*/
func MemoryAllMessage(args ...interface{}) ([]HistoryStruct, error) {
	con, ok := args[0].(*Conversation)
	if !ok {
		return nil, fmt.Errorf(memerr)
	}

	mem := []HistoryStruct{{Role: "system", Content: sp}}

	return append(mem, con.History...), nil
}

/*
MemoryLastMessage is a function that will use the last int x messages as memory.

@param The number of last messages to use as memory.

@return Array of the last X amount of messages
*/
func MemoryLastMessage(args ...interface{}) ([]HistoryStruct, error) {
	con, ok := args[0].(*Conversation)
	if !ok {
		return nil, fmt.Errorf(memerr)
	}

	count, ok := memArgs[0].(int)
	if !ok {
		return nil, fmt.Errorf("memory: second argument needs to be an integer")
	}

	count-- // Subtract 1 to account for binary counting

	if count < 0 {
		count = -count
	}

	historyLength := len(con.History)

	if count > historyLength {
		count = historyLength
	}

	mem := []HistoryStruct{{Role: "system", Content: sp}}

	// Always returns the system message first then the other interactions
	return append(mem, con.History[historyLength-count:]...), nil
}

/*
MemoryTime is a function that will use the messages within the last int x minutes as memory.

@param The number of last minutes to use as memory.

@return Array of the last X amount of messages within the last Y minutes
*/
func MemoryTime(args ...interface{}) ([]HistoryStruct, error) {
	con, ok := args[0].(*Conversation)
	if !ok {
		return nil, fmt.Errorf(memerr)
	}

	age, ok := memArgs[0].(int)
	if !ok {
		return nil, fmt.Errorf("memory: second argument needs to be an integer")
	}

	if age < 0 {
		age = -age
	}

	mem := []HistoryStruct{{Role: "system", Content: sp}}

	for i := range len(con.History) {
		if time.Since(con.History[i].TimeStamp).Minutes() < float64(age) {
			mem = append(mem, con.History[i])
		}
	}

	return mem, nil
}
