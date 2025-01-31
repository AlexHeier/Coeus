package llm

import (
	"Coeus/llm/tool"
	"Coeus/provider"
	"strings"
)

var conversations []Conversation

// Struct for containing the individual conversations with the LLMs
type Conversation struct {
	MainPrompt string
	ToolsResp  interface{}
	History    []string
	Summary    string
	UserPrompt string
}

// Appends a prompt and section to the history within the conversation
func (c *Conversation) AppendHistory(user, llm string) {
	newHistory := "\n[USER:]" + user + "\n[LLM:]" + llm
	c.History = append(c.History, newHistory)
}

func (c *Conversation) Prompt(s string) (map[string]interface{}, error) {
	var toolDesc string

	if len(tool.Tools) > 0 {
		for _, t := range tool.Tools {
			toolDesc += t.Name + ": " + t.Desc + "\n"
		}
	}

	prefix := c.MainPrompt + "[BEGIN TOOLS] " + tool.ToolDefintion + toolDesc + "[END TOOLS]\n[BEGIN HISTORY]" + Memory(c) + "[END HISTORY]\n"
	res, err := provider.Send(prefix + s)
	if err != nil {
		return res, err
	}
	resString := res["response"].(string)

	// Check for if the response contains a summary and extract it
	if strings.Contains(resString, "[BEGIN SUMMARY]") {
		beginIndex := strings.Index(resString, "[BEGIN SUMMARY]")
		endIndex := strings.Index(resString, "[END SUMMARY]")
		if endIndex > beginIndex {
			c.Summary = strings.TrimSpace(resString[beginIndex:endIndex])
		}
	}

	// TODO: Implement a way to use the tools

	c.AppendHistory(s, resString)
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
		MainPrompt: Persona + "Anser in the language the user is using.\n",
	}

	conversations = append(conversations, newCon)
	return &conversations[len(conversations)-1]
}
