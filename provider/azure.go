package provider

import (
	"Coeus/llm/tool"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Azure Configuration struct
type AzureProviderStruct struct {
	Endpoint    string  // Azure API endpoint
	APIKey      string  // Azure API Key
	Temperature float64 // How free thinking the LLM should be. Lower equals more free. Can be between 0.1 and 1.0
	MaxTokens   int     // Max amount of tokens a response can use
}

// Struct used in sending requests to an Azure endpoint
type AzureRequest struct {
	Messages    []AzureMessage `json:"messages"`
	Tools       []AzureTool    `json:"tools"`
	MaxTokens   int            `json:"max_tokens"`
	Temperature float64        `json:"temperature"`
}

// Struct for containing the response from Azure
type AzureResponse struct {
	Choices []struct {
		ContentFilterResults map[string]interface{} `json:"content_filter_results"`
		FinishReason         string                 `json:"finish_reason"`
		LogProbs             string                 `json:"logprobs"`
		Message              AzureMessage           `json:"message"`
	} `json:"choices"`
	Model               string `json:"model"`
	PromptFilterResults []struct {
		PromptIndex          int `json:"prompt_index"`
		ContentFilterResults map[string]struct {
			Filtered bool   `json:"filtered"`
			Severity string `json:"severity"`
		} `json:"content_filter_results"`
	} `json:"prompt_filter_results"`
	Usage struct {
		CompletionTokens int `json:"completion_tokens"`
		PromptTokens     int `json:"prompt_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

type AzureMessage struct {
	Content   string                   `json:"content"`
	Refusal   string                   `json:"refusal"`
	Role      string                   `json:"role"`
	ToolCalls []map[string]interface{} `json:"tool_calls"`
}

type AzureTool struct {
	Type     string `json:"type"`
	Function struct {
		Name        string                 `json:"name"`
		Description string                 `json:"description"`
		Parameters  map[string]interface{} `json:"parameters"`
		Required    []string               `json:"required"`
	} `json:"function"`
}

// Creates a new Azure config and sets it as provider
/*
@param endpoint: String which contains the URL of Azure endpoint
@param apikey: Azure API key
@param temperature: Specifies the amount of freedom an LLM should have when answering
@param maxTokens: Specifies the max amount of tokens an LLM answer should use
*/
func Azure(endpoint, apikey string, temperature float64, maxTokens int) error {
	if endpoint == "" {
		return fmt.Errorf("no endpoint specified")
	}

	if apikey == "" {
		return fmt.Errorf("no apikey specified")
	}

	if temperature < 0.1 || temperature > 1.0 {
		return fmt.Errorf("llm temperature set outside acceptable bounds")
	}

	if maxTokens <= 0 {
		return fmt.Errorf("maxtokens must be bigger than 0")
	}

	Provider = AzureProviderStruct{
		Endpoint:    endpoint,
		APIKey:      apikey,
		Temperature: temperature,
		MaxTokens:   maxTokens,
	}

	return nil
}

func SendAzure(request RequestStruct) (ResponseStruct, error) {

	Config := Provider.(AzureProviderStruct)

	AzureReq := CreateAzureRequest(request)

	buf := new(bytes.Buffer)

	json.NewEncoder(buf).Encode(AzureReq)

	req, err := http.NewRequest(http.MethodPost, Config.Endpoint, buf)
	if err != nil {
		return ResponseStruct{}, err
	}

	req.Header.Add("api-key", Config.APIKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return ResponseStruct{}, err
	}
	defer res.Body.Close()

	resData, err := io.ReadAll(res.Body)
	if err != nil {
		return ResponseStruct{}, err
	}

	var azureRes AzureResponse

	err = json.Unmarshal(resData, &azureRes)
	if err != nil {
		return ResponseStruct{}, err
	}

	switch azureRes.Choices[0].FinishReason {
	case "tool_calls":
		messages, err := AzureRunTools(azureRes)
		if err != nil {
			return ResponseStruct{}, err
		}
		*request.History = append(*request.History, messages...)
		return SendAzure(request)
	default:
		return ResponseStruct{Response: azureRes.Choices[0].Message.Content}, nil
	}
}

func AzureRunTools(res AzureResponse) ([]HistoryStruct, error) {
	var messages []HistoryStruct

	for _, call := range res.Choices[0].Message.ToolCalls {
		fmt.Println(call)

		f := call["function"].(map[string]interface{})
		args := make(map[string]interface{})

		jsonArgs := f["arguments"]
		toolName := f["name"]

		err := json.Unmarshal([]byte(jsonArgs.(string)), &args)
		if err != nil {
			return messages, err
		}

		_ = toolName

		//t, err := tool.Find(toolCall.Name)
		//if err != nil {
		//	return []AzureMessage{}, err
		//}

	}

	return messages, nil
}

func CreateAzureRequest(request RequestStruct) AzureRequest {
	Config := Provider.(AzureProviderStruct)

	AzureReq := AzureRequest{
		Temperature: Config.Temperature,
		MaxTokens:   Config.MaxTokens,
	}

	for _, h := range *request.History {
		AzureReq.Messages = append(AzureReq.Messages, AzureMessage{Role: h.Role, Content: h.Content})
	}

	for _, t := range tool.Tools {
		AzureReq.Tools = append(AzureReq.Tools, AzureTool{Type: "function", Function: struct {
			Name        string                 "json:\"name\""
			Description string                 "json:\"description\""
			Parameters  map[string]interface{} "json:\"parameters\""
			Required    []string               "json:\"required\""
		}{Name: t.Name,
			Description: t.Desc,
			Parameters:  t.Params}})
	}

	AzureReq.Messages = append(AzureReq.Messages, AzureMessage{Role: "system", Content: request.Systemprompt}, AzureMessage{Role: "user", Content: request.Userprompt})

	return AzureReq
}
