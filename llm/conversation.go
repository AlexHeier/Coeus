package llm

import "Coeus/provider"

var conversations []Conversation

// Struct for containing the individual conversations with the LLMs
type Conversation struct {
	MainPrompt     string
	ToolsResp      interface{}
	History        []string
	Summary        string
	LatestResponse string
	UserPrompt     string
}

// Appends a prompt and section to the history within the conversation
func (c *Conversation) AppendHistory(s, section string) {
	newHistory := "[" + section + "]\n" + s + "\n\n"
	c.History = append(c.History, newHistory)
}

func (c *Conversation) Prompt(s string) (map[string]interface{}, error) {
	prefix := c.MainPrompt + "[BEGIN HISTORY]" + c.DumpConversation() + "[END HISTORY]\n\n"
	res, err := provider.Send(prefix + s)
	if err != nil {
		return res, err
	}
	c.AppendHistory(s, "USER")
	c.AppendHistory(res["response"].(string), "LLM")
	c.LatestResponse = res["response"].(string)
	return res, err
}

func (c *Conversation) DumpConversation() string {
	dump := c.MainPrompt

	for _, h := range c.History {
		dump += h
	}
	return dump
}

func BeginConversation() *Conversation {
	newCon := Conversation{
		MainPrompt: Persona_HulkHogan,
	}

	conversations = append(conversations, newCon)
	return &conversations[len(conversations)-1]
}
