package llm

var Message Struct

type Struct struct {
	llm    interface{}
	Create func(interface{}) (string, error)
}
