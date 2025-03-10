package coeus

import (
	"fmt"
)

/* Default memory function is All. This function will use all messages as memory. */
var memory func(args ...interface{}) string = MemoryAllMessage
var memArgs []interface{}

/*
MemoryVersion changes the function used for memory management. Default is All messages.

@param newFunc: The new function to use for memory management.
*/
func MemoryVersion(newFunc ...interface{}) {
	if len(newFunc) > 0 {
		if fn, ok := newFunc[0].(func(args ...interface{}) string); ok {
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

@return A string of the history.
*/
func MemoryAllMessage(args ...interface{}) string {
	con, ok := args[0].(*Conversation)
	if !ok {
		fmt.Println("MEMORY: Bad type. How?")
		return ""
	}

	if len(con.History) <= 0 {
		return "No History"
	}

	var temp string
	for _, h := range con.History {
		temp += h.Role + ": " + h.Content
	}

	return temp
}

/*
MemoryLastMessage is a function that will use the last int x messages as memory.

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
	var hist string
	for _, h := range lastMessages {
		hist += h.Role + ": " + h.Content + "\n"
	}

	return "[BEGIN HISTORY]" + hist + "[END HISTORY]"
}

/*
MemorySummary will create a summary of the conversation and use it as memory.

@return A string representing the new message with the summary and an error if the conversion fails.
*/
func MemorySummary(args ...interface{}) string {
	var tempSummary string
	makeNewSummary := `Respond to the user's latest message appropriately. Then, generate a summary using only the previous summary [OLD SUMMARY] and the latest user message. Do not include any other context. Begin the summary with "[SUMMARY]". Ensure that all information from the old summary is retained.`

	con, ok := args[0].(*Conversation)
	if !ok {
		return "No history"
	}

	if con.Summary == "" {
		for _, h := range con.History {
			tempSummary += h.Role + ": " + h.Content + "\n"
		}
		con.Summary = tempSummary
	}

	return "[BEGIN OLD SUMMARY] " + con.Summary + " [END OLD SUMMARY]" + makeNewSummary
}
