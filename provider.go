package coeus

import "fmt"

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
}

type ToolCall struct {
	Index    *int   `json:"index,omitempty"`
	ID       string `json:"id,omitempty"`
	Type     string `json:"type"`
	Function struct {
		Name      string                 `json:"name,omitempty"`
		Arguments map[string]interface{} `json:"arguments,omitempty"`
	} `json:"function"`
}

// Provider will change to the provider struct of the chosen provider
var Provider interface{}

/*
Send function will send the request to any provider and return the response

@param request RequestStruct

@return ResponseStruct, error
*/
func Send(request RequestStruct) (ResponseStruct, error) {

	switch Provider.(type) {
	case ollamaStruct:
		return sendOllama(request)
	case azureProviderStruct:
		return sendAzure(request)
	case openAIStruct:
		return sendOpenAI(request)
	default:
		return ResponseStruct{}, fmt.Errorf("no valid provider configured")
	}
}
