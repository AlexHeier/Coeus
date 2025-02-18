package provider

import "fmt"

type ResponseStruct struct {
	Response             string
	LoadDuration         float64
	eval_count           float64
	prompt_eval_count    float64
	prompt_eval_duration float64
}

type RequestStruct struct {
	History      []HistoryStruct
	Systemprompt string
	Userprompt   string
}

type HistoryStruct struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

var Provider interface{}

func Send(request RequestStruct) (ResponseStruct, error) {

	switch Provider.(type) {
	case OllamaStruct:
		return SendOllama(request)
	case AzureStruct:
		return SendAzure(request)
	case OpenAIStruct:
		return SendOpenAI(request)
	default:
		return ResponseStruct{}, fmt.Errorf("no valid provider configured")
	}
}
