package memory

/*
Summary is a function that will take a summary of the conversation and use the summary as memory.

@param intf An interface{} that should be a slice of strings representing the conversation.
@return A string representing the new message with the summary and an error if the conversion fails.
*/

// TODO: Finne ut av hvordan man skal lage denne
func Summary(args ...interface{}) string {
	return SummaryString + "\nCreate a summary of the new conversation and add it as Summary:"
}
