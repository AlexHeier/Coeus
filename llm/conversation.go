package llm

import "Coeus/provider"

var conversations []Conversation

type Conversation struct {
	MainPrompt     string
	ToolsDesc      string
	ToolsResp      string
	History        []string
	LatestResponse string
	UserPrompt     string
}

func (c *Conversation) AppendUserHistory(s string) {
	newHistory := "[USER]\n" + s + "\n\n"
	c.History = append(c.History, newHistory)
}

func (c *Conversation) AppendSystemHistory(s string) {
	newHistory := "[SYSTEM]\n" + s + "\n\n"
	c.History = append(c.History, newHistory)
}

func (c *Conversation) AppendLLMHistory(s string) {
	newHistory := "[LLM]\n" + s + "\n\n"
	c.History = append(c.History, newHistory)
}

func (c *Conversation) Prompt(s string) (map[string]interface{}, error) {
	prefix := c.MainPrompt + "[BEGIN HISTORY]" + c.DumpConversation() + "[END HISTORY]\n\n"
	res, err := provider.Send(prefix + s)
	if err != nil {
		return res, err
	}
	c.AppendUserHistory(s)
	c.AppendLLMHistory(res["response"].(string))
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
		MainPrompt: Persona_BarackObama,
	}

	conversations = append(conversations, newCon)
	return &conversations[len(conversations)-1]
}
