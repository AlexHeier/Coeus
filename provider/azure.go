package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/AlexHeier/Coeus/llm/tool"
)

const AZURE_ROLE_USER = "user"
const AZURE_ROLE_TOOL = "tool"
const AZURE_ROLE_SYSTEM = "system"
const AZURE_ROLE_ASSISTANT = "assistant"
const AZURE_TYPE_FUNCTION = "function"

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

func sendAzure(request RequestStruct) (ResponseStruct, error) {

	azureRes, err := azureSendRequest(createAzureRequest(request))
	if err != nil {
		return ResponseStruct{}, err
	}

	if len(azureRes.Choices[0].Message.ToolCalls) > 0 {

		// Push LLM tool calls to the history
		*request.History = append(*request.History, HistoryStruct{Role: AZURE_ROLE_ASSISTANT,
			ToolCalls: azureRes.Choices[0].Message.ToolCalls})

		for _, toolCall := range azureRes.Choices[0].Message.ToolCalls {

			t, err := tool.Find(toolCall.Function.Name)
			if err != nil {
				return ResponseStruct{}, fmt.Errorf("could not find the tool %s", t.Name)
			}

			parsedToolCall := make(map[string]interface{})
			err = json.Unmarshal([]byte(toolCall.Function.Arguments), &parsedToolCall)
			if err != nil {
				return ResponseStruct{}, fmt.Errorf("failed to parse tool arguments: %v", err)
			}

			var args []interface{}
			for _, val := range parsedToolCall {
				args = append(args, val)
			}

			toolResponse, err := t.Run(args...)
			if err != nil {
				return ResponseStruct{}, fmt.Errorf("error during tool execution: %v", err)
			}

			fmt.Println(toolResponse)

			*request.History = append(*request.History, HistoryStruct{
				Role:       AZURE_ROLE_TOOL,
				Content:    toolResponse,
				ToolCallID: toolCall.ID})
		}

		azureRes, err = azureSendRequest(createAzureRequest(request))
		if err != nil {
			return ResponseStruct{}, err
		}
	}
	return ResponseStruct{Response: azureRes.Choices[0].Message.Content}, nil
}

func createAzureRequest(request RequestStruct) azureRequest {
	Config := Provider.(AzureProviderStruct)

	AzureReq := azureRequest{
		Temperature: Config.Temperature,
		MaxTokens:   Config.MaxTokens,
	}

	AzureReq.Messages = append(AzureReq.Messages, azureMessage{
		Role:    AZURE_ROLE_SYSTEM,
		Content: request.Systemprompt})

	*request.History = append(*request.History, HistoryStruct{
		Role:    AZURE_ROLE_USER,
		Content: request.Userprompt})

	for _, h := range *request.History {
		AzureReq.Messages = append(AzureReq.Messages,
			azureMessage{
				Role:       h.Role,
				Content:    h.Content,
				ToolCallID: h.ToolCallID,
				ToolCalls:  h.ToolCalls})
	}

	for _, t := range tool.Tools {
		AzureReq.Tools = append(AzureReq.Tools, azureTool{Type: AZURE_TYPE_FUNCTION, Function: struct {
			Name        string   `json:"name"`
			Description string   `json:"description"`
			Parameters  any      `json:"parameters"`
			Required    []string `json:"required"`
		}{Name: t.Name,
			Description: t.Desc,
			Parameters:  t.Params}})
	}

	return AzureReq
}

func azureSendRequest(AzureReq azureRequest) (azureResponse, error) {
	Config := Provider.(AzureProviderStruct)

	buf := new(bytes.Buffer)

	json.NewEncoder(buf).Encode(AzureReq)

	req, err := http.NewRequest(http.MethodPost, Config.Endpoint, buf)
	if err != nil {
		return azureResponse{}, err
	}

	req.Header.Add("api-key", Config.APIKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return azureResponse{}, err
	}
	defer res.Body.Close()

	err = azureParseStatusCode(res)
	if err != nil {
		return azureResponse{}, err
	}

	resData, err := io.ReadAll(res.Body)
	if err != nil {
		return azureResponse{}, err
	}

	var azureRes azureResponse

	err = json.Unmarshal(resData, &azureRes)
	if err != nil {
		return azureResponse{}, err
	}

	return azureRes, nil
}

func azureParseStatusCode(res *http.Response) error {
	switch res.StatusCode {
	case 429: // Token limit reached
		return fmt.Errorf("token rate limit reached")
	case 200: // OK
		return nil
	default:
		fmt.Printf("unhandled status code: %d\n", res.StatusCode)
		return nil
	}
}
