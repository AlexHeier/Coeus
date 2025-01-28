package llm

func Setup(llm interface{}, memoryFunc func(interface{}) (string, error)) {
	Message.llm = llm
	Message.Create = memoryFunc
}
