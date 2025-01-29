package llm

var MainConversation Conversation

type Conversation struct {
	MainPrompt   string
	ToolsDesc    string
	ToolRepsonse string
	History      []string
	UserPrompt   string
}

func Setup(llm interface{}, memoryFunc func(interface{}) (string, error)) {
	Message.llm = llm
	Message.Create = memoryFunc
}
