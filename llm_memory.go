package coeus

import (
	"fmt"
)

/* Default memory function is All. This function will use all messages as memory. */
var memory func(args ...interface{}) ([]HistoryStruct, error) = MemoryAllMessage
var memArgs []interface{}

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
		return nil, fmt.Errorf("MEMORY: Bad type. How?")
	}

	fmt.Println("Hello")

	return append(con.History, HistoryStruct{Role: "system", Content: sp}), nil
}

/*
MemoryLastMessage is a function that will use the last int x messages as memory.

@param The number of last messages to use as memory.

@return Array of the last X amount of messages
*/
func MemoryLastMessage(args ...interface{}) ([]HistoryStruct, error) {
	con, ok := args[0].(*Conversation)
	if !ok {
		return nil, fmt.Errorf("BAD CONVERSATION: How?")
	}

	elements, ok := memArgs[0].(int)
	if !ok {
		return nil, fmt.Errorf("second argument needs to be an integer")
	}

	if elements < 0 {
		return nil, fmt.Errorf("integer needs to be a positive number")
	}

	historyLen := len(con.History)
	fmt.Println(historyLen)
	if elements > historyLen {
		elements = historyLen
	}

	system := []HistoryStruct{{Role: "system", Content: sp}}

	// Always returns the system message first then the other interactions
	return append(system, con.History[:historyLen-elements]...), nil
}

/*
CURRENTLY NOT IN USE
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
