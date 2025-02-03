package provider

import "fmt"

type ResponseStruct struct {
	Response             string
	LoadDuration         float64
	eval_count           float64
	prompt_eval_count    float64
	prompt_eval_duration float64
}

var Provider interface{}

func Send(prompt string) (ResponseStruct, error) {

	switch Provider.(type) {
	case OllamaStruct:
		return SendOllama(prompt)
	case AzureStruct:
		return SendAzure(prompt)
	case OpenAIStruct:
		return SendOpenAI(prompt)
	default:
		return ResponseStruct{}, fmt.Errorf("no valid provider configured")
	}
}
