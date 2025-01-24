package conversation

type memoryFunc func([]string) []string

type llmInput interface{}

type llmFunc func(llmInput) struct{}

type ConversationSetup struct {
	llm    llmFunc
	Memory memoryFunc
}
