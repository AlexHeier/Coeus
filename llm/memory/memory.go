package memory

// Dette er bare eksempler, vi må lage dem bedre senere. Tror vi bør sette opp pekere her, tanken er at vi definerer memory=func der func er en av de neden for.
// Så bruker programet riktig minnehåndtering deretter. Evt så må man bare gjøre et systemkall typ newMessage = Coeus.Conversation.Memory.Summery(OldSummery, UserMessage) eller tilsvarende.
// evt sette opp en funksjon som bytter imellom All() og Summery basert på billigst token bruk fortløpende.

// alle må ha en interface{} som input og returnere en string og en error.

/**
* Last is a function that will use the last int x messages as memory
 */
func Last(int) {

}
