package memory

/*
All is a function that will use all messages as memory.

@return A string of the history.
*/
func All(args ...interface{}) string {
	var temp string
	for _, h := range args.History {
		temp += h
	}
	return temp
}
