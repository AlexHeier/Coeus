package memory

import "fmt"

// Dette er bare eksempler, vi må lage dem bedre senere. Tror vi bør sette opp pekere her, tanken er at vi definerer memory=func der func er en av de neden for.
// Så bruker programet riktig minnehåndtering deretter. Evt så må man bare gjøre et systemkall typ newMessage = Coeus.Conversation.Memory.Summery(OldSummery, UserMessage) eller tilsvarende.
// evt sette opp en funksjon som bytter imellom All() og Summery basert på billigst token bruk fortløpende.

/**
* Summery is a function that will take a summery of the conversation and use the summary as memory.
 */
func Summary(chat []string) []string {
	prompt := `You are an AI and will help answer questions from a human. 
	You will anwser to the best of your ability. If you do not know the answer, you will say so. 
	Answer in the language the user is using. 
	Summery is the summery of the conversation, human is the message form the human. 
	Anser with AI: {x} where x is the answer to the users request and NewSummery {x} which is a new summery of the conversation.`

	newMessage := []string{}
	newMessage = append(newMessage, fmt.Sprintf("%v\n Summary: {%v}\nHuman: {\n}", prompt, chat))
	return newMessage
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
