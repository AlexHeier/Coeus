package coeus

// sp is the system prompt the LLM will use
var sp string

/*
SetSystemPrompt is a function that sets the system prompt for the LLM.

@param prompt: the system prompt to set
*/
func SetSystemPrompt(prompt string) {
	sp = prompt + "\n"
}
