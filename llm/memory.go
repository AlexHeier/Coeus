package llm

import (
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
	con, ok := args[0].(*Conversation)
	if !ok {
		return "No history"
	}

	var temp string
	for _, h := range con.History {
		temp += h
	}

	print(temp)

	return "[BEGIN HISTORY]" + temp + "[END HISTORY]"
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
	return "[BEGIN HISTORY]" + strings.Join(lastMessages, "\n") + "[END HISTORY]"
}

/*
Summary will create a summary of the conversation and use it as memory.

@return A string representing the new message with the summary and an error if the conversion fails.
*/
func MemorySummary(args ...interface{}) string {
	var tempSum string

	con, ok := args[0].(*Conversation)
	if !ok {
		return "No history"
	}

	if con.Summary == "" {
		for _, h := range con.History {
			tempSum += h + "\n"
		}
		tempSum += "Can you create a short summary of the conversation that contains everything important at the end of the message? begin the summary with [BEGIN SUMMARY] and end with [END SUMMARY]."
		con.Summary = tempSum
	} else {
		con.Summary = "[BEGIN OLD SUMMARY]" + con.Summary + "[END OLD SUMMARY]\nCreate a precise summary of the last summary + the users new message. Add the new summary at the end of the message? begin the summary with [BEGIN SUMMARY] and end with [END SUMMARY]."
	}

	return con.Summary
}
