package llm

import (
	"Coeus/provider"
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
It's inefficient as it takes an additional LLM call to create the history summary. Other options should be considered.

@param intf An interface{} that should be a slice of strings representing the conversation.
@return A string representing the new message with the summary and an error if the conversion fails.
*/
func MemorySummary(args ...interface{}) string {

	con, ok := args[0].(*Conversation)
	if !ok {
		return "No history"
	}

	messages, ok := args[1].(int)
	if !ok {
		return "Amount of history to include not specified as an int"
	}

	// Makes sure that the amount of messages remains a positive number
	if (messages < 0) || (messages > 10) {
		fmt.Printf("MemorySummary: Amount of messages outside bounds! Using internal value of %d. Was %d\n", 2, messages)
		messages = 2
	}

	// Summary is set to contain atleast 4 entries of history. That way the AI has more information to work with and gives better answers.
	if len(con.History) < (messages + 4) {
		fmt.Printf("Not enough history to create summary. Need %d or more. Have %d\n", messages+4, len(con.History))
		var dump string
		for _, history := range con.History {
			dump += history
		}
		return dump
	}

	var sumPrompt string
	for _, history := range con.History[:len(con.History)-messages] {
		sumPrompt += history
	}

	res, err := provider.Send("Only create a short bulletpoint summary of this conversation between the user and LLM. Highlights things which will help the LLM to keep the conversation going. No other text.\n[BEGIN]\n" + sumPrompt + "\n[END]\n")
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}

	con.Summary = "[BEGIN SUMMARY]\n" + res.Response + "\n[END SUMMARY]\n"

	var ret string
	for _, h := range con.History[messages:] {
		ret += h
	}

	return con.Summary + ret
}
