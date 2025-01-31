package memory

import (
	"Coeus/llm"
	"strings"
)

/*
Last is a function that will use the last int x messages as memory.

@param The number of last messages to use as memory.
@return A string representing the memory.
*/
func Last(args ...interface{}) string {
	if len(args) == 0 {
		return ""
	}

	x, ok := args[0].(int)
	if !ok {
		return ""
	}

	historyLen := len(llm.History)
	if x > historyLen {
		x = historyLen
	}

	lastMessages := llm.History[historyLen-x:]
	return strings.Join(lastMessages, "\n")
}
