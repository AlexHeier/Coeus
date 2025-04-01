package coeus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const azureRoleUser = "user"
const azureRoleTool = "tool"

// const azureRoleSystem = "system"
const azureRoleAssistant = "assistant"
const azureTypeFunction = "function"

/*
Azure sets the provider to Azure.

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

	Provider = azureProviderStruct{
		Endpoint:    endpoint,
		APIKey:      apikey,
		Temperature: temperature,
		MaxTokens:   maxTokens,
	}

	return nil
}

/*
sendAzure sends a request to Azure.

@param request: the request to send

@return A response and an error if the request fails
*/
func sendAzure(con *Conversation) (ResponseStruct, error) {

	azureRes, err := azureSendRequest(createAzureRequest(con))
	if err != nil {
		return ResponseStruct{}, err
	}

	if len(azureRes.Choices[0].Message.ToolCalls) > 0 {

		for {

			con.History = append(con.History, HistoryStruct{
				Role:      azureRoleAssistant,
				ToolCalls: azureRes.Choices[0].Message.ToolCalls,
			})

			for _, toolCall := range azureRes.Choices[0].Message.ToolCalls {

				t, err := FindTool(toolCall.Function.Name)
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

				toolResponse, err := t.RunTool(args...)
				if err != nil {
					return ResponseStruct{}, fmt.Errorf("error during tool execution: %v", err)
				}

				con.History = append(con.History, HistoryStruct{
					Role:       azureRoleTool,
					Content:    toolResponse,
					ToolCallID: toolCall.ID,
				})

			}

			azureRes, err = azureSendRequest(createAzureRequest(con))
			if err != nil {
				return ResponseStruct{}, err
			}

			if len(azureRes.Choices[0].Message.ToolCalls) == 0 {
				break
			}
		}
	}

	return ResponseStruct{Response: azureRes.Choices[0].Message.Content}, nil
}

/*
createAzureRequest creates an Azure request.

@param request: the request to send

@return: an azureRequest struct
*/
func createAzureRequest(con *Conversation) azureRequest {
	Config := Provider.(azureProviderStruct)

	AzureReq := azureRequest{
		Temperature: Config.Temperature,
		MaxTokens:   Config.MaxTokens,
	}

	history, err := memory(con)
	if err != nil {
		history = []HistoryStruct{{Role: "system", Content: sp}}
	}

	for _, h := range history {
		AzureReq.Messages = append(AzureReq.Messages,
			azureMessage{
				Role:       h.Role,
				Content:    h.Content,
				ToolCallID: h.ToolCallID,
				ToolCalls:  h.ToolCalls})
	}

	if rag {
		AzureReq.Messages = append(AzureReq.Messages,
			azureMessage{
				Role:    "system",
				Content: GetRAG(con.History[len(con.History)-1].Content)})
	}

	for _, t := range Tools {
		AzureReq.Tools = append(AzureReq.Tools, azureTool{Type: azureTypeFunction, Function: struct {
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

/*
azureSendRequest sends a request to Azure.

@param AzureReq: the request to send

@return: an azureResponse struct and an error if the request fails
*/
func azureSendRequest(AzureReq azureRequest) (azureCompletionMessage, error) {
	Config := Provider.(azureProviderStruct)

	buf := new(bytes.Buffer)

	json.NewEncoder(buf).Encode(AzureReq)

	req, err := http.NewRequest(http.MethodPost, Config.Endpoint, buf)
	if err != nil {
		return azureCompletionMessage{}, err
	}

	req.Header.Add("api-key", Config.APIKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return azureCompletionMessage{}, err
	}
	defer res.Body.Close()

	err = azureParseStatusCode(res)
	if err != nil {
		return azureCompletionMessage{}, err
	}

	resData, err := io.ReadAll(res.Body)
	if err != nil {
		return azureCompletionMessage{}, err
	}

	var azureRes azureCompletionMessage

	err = json.Unmarshal(resData, &azureRes)
	if err != nil {
		return azureCompletionMessage{}, err
	}

	return azureRes, nil
}

/*
azureParseStatusCode parses the status code of an Azure response.

@param res: the response to parse

@return: an error if the status code is not 200
*/
func azureParseStatusCode(res *http.Response) error {
	switch res.StatusCode {
	case 429: // Token limit reached
		return fmt.Errorf("token rate limit reached")
	case 200: // OK
		return nil
	default:
		data, _ := io.ReadAll(res.Body)
		return fmt.Errorf("azure: unhandled status code: %d, response: %s", res.StatusCode, string(data))
	}
}
