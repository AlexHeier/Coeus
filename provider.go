package coeus

import (
	"fmt"
	"time"
)

type ResponseStruct struct {
	Response           string
	TotalLoadDuration  float64
	Eval_count         float64
	PromptEvalCount    float64
	PromptEvalDuration float64
}

type RequestStruct struct {
	History      *[]HistoryStruct
	Systemprompt string
	Userprompt   string
}

type HistoryStruct struct {
	Role       string     `json:"role"`
	Content    string     `json:"content"`
	ToolCallID string     `json:"tool_call_id,omitempty"`
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
	TimeStamp  time.Time
}

type ToolCall struct {
	Index    *int   `json:"index,omitempty"`
	ID       string `json:"id,omitempty"`
	Type     string `json:"type"`
	Function struct {
		Name      string `json:"name,omitempty"`
		Arguments string `json:"arguments,omitempty"`
	} `json:"function"`
}

// Provider will change to the provider struct of the chosen provider
var Provider interface{}

/*
Send function will send the request to any provider and return the response

@param con ConversationStruct

@return ResponseStruct, error
*/
func Send(con *Conversation) (ResponseStruct, error) {

	switch Provider.(type) {
	case ollamaStruct:
		return sendOllama(con)
	case azureProviderStruct:
		return sendAzure(con)
	case openAIStruct:
		return sendOpenAI(con)
	default:
		return ResponseStruct{}, fmt.Errorf("no valid provider configured")
	}
}
