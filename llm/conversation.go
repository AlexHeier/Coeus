package llm

type Conversation struct {
	MainPrompt string
	ToolsDesc  string
	ToolsResp  string
	History    string
	UserPrompt string
}

var MainConversation Conversation

func (c *Conversation) UpdateHistory(s string) {
	c.History += s
}
