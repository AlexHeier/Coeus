package provider

import (
	"Coeus/llm/tool"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

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

	azureRes, err := AzureSendRequest(CreateAzureRequest(request))
	if err != nil {
		return ResponseStruct{}, err
	}

	if len(azureRes.Choices[0].Message.ToolCalls) > 0 {
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

			*request.History = append(*request.History, HistoryStruct{Role: "tool", Content: toolResponse, ToolCallID: toolCall.ID})

		}

		fmt.Println(CreateAzureRequest(request))

		azureRes, err = AzureSendRequest(CreateAzureRequest(request))
		if err != nil {
			return ResponseStruct{}, err
		}

		fmt.Println(azureRes)

	}

	return ResponseStruct{Response: azureRes.Choices[0].Message.Content}, nil
}

func CreateAzureRequest(request RequestStruct) azureRequest {
	Config := Provider.(AzureProviderStruct)

	AzureReq := azureRequest{
		Temperature: Config.Temperature,
		MaxTokens:   Config.MaxTokens,
	}

	for _, h := range *request.History {
		AzureReq.Messages = append(AzureReq.Messages, azureMessage{Role: h.Role, Content: h.Content})
	}

	for _, t := range tool.Tools {
		AzureReq.Tools = append(AzureReq.Tools, azureTool{Type: "function", Function: struct {
			Name        string   "json:\"name\""
			Description string   "json:\"description\""
			Parameters  any      "json:\"parameters\""
			Required    []string "json:\"required\""
		}{Name: t.Name,
			Description: t.Desc,
			Parameters:  t.Params}})
	}

	AzureReq.Messages = append(AzureReq.Messages, azureMessage{Role: "system", Content: request.Systemprompt}, azureMessage{Role: "user", Content: request.Userprompt})

	return AzureReq
}

func AzureSendRequest(AzureReq azureRequest) (azureResponse, error) {
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

	resData, err := io.ReadAll(res.Body)
	if err != nil {
		return azureResponse{}, err
	}

	fmt.Println(string(resData))

	var azureRes azureResponse

	err = json.Unmarshal(resData, &azureRes)
	if err != nil {
		return azureResponse{}, err
	}

	return azureRes, nil
}
