package memory

/*
All is a function that will use all messages as memory.

@param h is the history of messages.
@return A string of the history.
*/
func All(h []string) string {
	var temp string
	for _, v := range h {
		temp += v + "\n"
	}
	return temp
}
