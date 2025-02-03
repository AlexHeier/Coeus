package llm

import (
	"fmt"
	"strings"
)

/* Default memory function is All. This function will use all messages as memory. */
var Memory func(args ...interface{}) string = MemoryAllMessage
var MemArgs []interface{}

/* Changes the function used for memory managment. Default is All messages. */
func MemoryVersion(newFunc ...interface{}) {
	if len(newFunc) > 0 {
		if fn, ok := newFunc[0].(func(args ...interface{}) string); ok {
			Memory = fn
			if len(newFunc) > 1 {
				MemArgs = newFunc[1:]
			}
		}
	} else {
		Memory = MemoryAllMessage // Default
		return
	}
}

/*
All is a function that will use all messages as memory.

@return A string of the history.
*/
func MemoryAllMessage(args ...interface{}) string {
	fmt.Println(args[0])
	con, ok := args[0].(*Conversation)
	if !ok {
		return "No history"
	}

	var temp string
	for _, h := range con.History {
		temp += h
	}
	return temp
}

/*
Last is a function that will use the last int x messages as memory.

@param The number of last messages to use as memory.
@return A string representing the memory.
*/
func MemoryLastMessage(args ...interface{}) string {
	con, ok := args[0].(*Conversation)
	if !ok {
		return "No history"
	}

	x, ok := args[1].(int)
	if !ok {
		return ""
	}

	historyLen := len(con.History)
	if x > historyLen {
		x = historyLen
	}

	lastMessages := con.History[historyLen-x:]
	return strings.Join(lastMessages, "\n")
}

/*
Summary is a function that will take a summary of the conversation and use the summary as memory.

@param intf An interface{} that should be a slice of strings representing the conversation.
@return A string representing the new message with the summary and an error if the conversion fails.
*/

func MemorySummary(args ...interface{}) string {
	con, ok := args[0].(*Conversation)
	if !ok {
		return "No history"
	}

	if con.Summary == "" {
		for _, h := range con.History {
			con.Summary += h
		}
	}

	return con.Summary + "\nCreate a new summary of the conversation at the end of the conversation. Start with [BEGIN SUMMARY] and end with [END SUMMARY]."
}
