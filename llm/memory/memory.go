package memory

/*
Default memory function is All. This function will use all messages as memory.
*/
var Memory interface{} = All

/*
Changes the function used for memory managment. Default is All messages.
*/
func Version(newFunc interface{}) {
	Memory = newFunc
}
