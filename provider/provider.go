package provider

import "fmt"

var Provider interface{}

func Send(prompt string) (map[string]interface{}, error) {

	switch Provider.(type) {
	case OllamaStruct:
		return SendOllama(prompt)
	case AzureStruct:
		return SendAzure(prompt)
	default:
		var v map[string]interface{}
		return v, fmt.Errorf("wrong provider type")
	}
}
