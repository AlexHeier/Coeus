package memory

import "fmt"

// Dette er bare eksempler, vi må lage dem bedre senere

/**
* Summery is a function that will take a summery of the conversation and use the summary as memory.
 */

func Summery(OldSummery, UserMessage string) (NewMessage string) {
	prompt := `You are an AI and will help answer questions from a human. 
	You will anwser to the best of your ability. If you do not know the answer, you will say so. 
	Answer in the language the user is using. 
	Summery is the summery of the conversation, human is the message form the human. 
	Anser with AI: {x} where x is the answer to the users request and NewSummery which is a new summery of the conversation.`

	NewMessage = fmt.Sprintf("%v\n Summery: {%v}\nHuman: {%v\n}", prompt, OldSummery, UserMessage)
	return NewMessage
}

/**
* Last is a function that will use the last int x messages as memory
 */
func Last(int) {

}

/**
* All is a function that will use all messages as memory.
 */
func All() {

}
