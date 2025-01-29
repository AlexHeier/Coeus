package memory

/**
* All is a function that will use all messages as memory.
 */
func All(a []string) string {
	var temp string
	for _, v := range a {
		temp += v + "\n"
	}
	return temp
}
