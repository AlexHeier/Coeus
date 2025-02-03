package provider

import "fmt"

type ResponseStruct struct {
	Response             string
	LoadDuration         string
	eval_count           string
	prompt_eval_count    string
	prompt_eval_duration string
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
