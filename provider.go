package coeus

import "fmt"

type ResponseStruct struct {
	Response             string
	TotalLoadDuration    float64
	Eval_count           float64
	Prompt_eval_count    float64
	Prompt_eval_duration float64
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

var Provider interface{}
var TTSProvider interface{}

func Send(request RequestStruct) (ResponseStruct, error) {

	switch Provider.(type) {
	case OllamaStruct:
		return SendOllama(request)
	case AzureProviderStruct:
		return sendAzure(request)
	case OpenAIStruct:
		return SendOpenAI(request)
	default:
		return ResponseStruct{}, fmt.Errorf("no valid provider configured")
	}
}
