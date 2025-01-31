package memory

/* Default memory function is All. This function will use all messages as memory. */
var Memory func(args ...interface{}) string = All
var Args []interface{}

/* The summary of the conversation, if the user uses summary memory. */
var SummaryString string

/* Changes the function used for memory managment. Default is All messages. */
func Version(newFunc ...interface{}) {
	if len(newFunc) > 0 {
		if fn, ok := newFunc[0].(func(args ...interface{}) string); ok {
			Memory = fn
			if len(newFunc) > 1 {
				Args = newFunc[1:]
			}
		}
	} else {
		Memory = All // Default
		return
	}
}
