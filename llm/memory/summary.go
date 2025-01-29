package memory

import (
	"fmt"
	"strings"
)

/**
* Summery is a function that will take a summery of the conversation and use the summary as memory.
 */
func Summary(intf interface{}) (string, error) {

	chat, ok := intf.([]string)
	if ok == false {
		return "", fmt.Errorf("could not convert interface to []string")
	}

	prompt := `You are an AI and will help answer questions from a human. 
	You will anwser to the best of your ability. If you do not know the answer, you will say so. 
	Answer in the language the user is using. 
	Summery is the summery of the conversation, human is the message form the human. 
	Anser with AI: {x} where x is the answer to the users request and NewSummery {x} which is a new summery of the conversation.`

	newMessage := prompt + "\n" + strings.Join(chat, "\n")
	return newMessage, nil
}
