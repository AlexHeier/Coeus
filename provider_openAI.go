package coeus

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

/*
OpenAI is a function that sets the provider to OpenAI.

@param model: the model to use

@param apikey: the api key to use

@return An error if the model or api key is empty
*/
func OpenAI(model, apikey string) error {
	Provider = openAIStruct{
		Model:  model,
		ApiKey: apikey,
	}
	return nil
}

/*
sendOpenAI is a function that sends a request to OpenAI.

@param request: the request to send

@return A response and an error if the request fails
*/
func sendOpenAI(con *Conversation) (ResponseStruct, error) {
	config := Provider.(openAIStruct)
	client := openai.NewClient(config.ApiKey)

	openAITools := convertToOpenAITools()

	req := openai.ChatCompletionRequest{
		Model:    config.Model,
		Tools:    openAITools,
		Messages: createOpenAIMessages(con),
	}

	resp, err := client.CreateChatCompletion(context.TODO(), req)
	if err != nil {
		return ResponseStruct{}, err
	}

	if len(resp.Choices[0].Message.ToolCalls) > 0 {

		for {

			con.History = append(con.History, HistoryStruct{
				Role:      "assistant",
				ToolCalls: convertToHistoryToolCalls(resp.Choices[0].Message.ToolCalls),
			})

			for _, t := range resp.Choices[0].Message.ToolCalls {
				tool, err := FindTool(t.Function.Name)
				if err != nil {
					return ResponseStruct{}, err
				}

				var parsedToolCall = make(map[string]interface{})
				err = json.Unmarshal([]byte(t.Function.Arguments), &parsedToolCall)
				if err != nil {
					return ResponseStruct{}, fmt.Errorf("failed to parse tool arguments: %v", err)
				}

				var args []interface{}
				for _, val := range parsedToolCall {
					args = append(args, val)
				}

				toolResponse, err := tool.RunTool(args...)
				if err != nil {
					return ResponseStruct{}, err
				}

				con.History = append(con.History, HistoryStruct{
					Role:       "tool",
					ToolCallID: t.ID,
					Content:    toolResponse,
				})
			}
			resp, err = client.CreateChatCompletion(context.TODO(), openai.ChatCompletionRequest{
				Model:    config.Model,
				Messages: createOpenAIMessages(con),
				Tools:    openAITools,
			})
			if err != nil {
				return ResponseStruct{}, err
			}

			if len(resp.Choices[0].Message.ToolCalls) == 0 {
				break
			}
		}
	}
	return ResponseStruct{Response: resp.Choices[0].Message.Content}, nil
}

/*
convertToOpenAITools is a function that converts the tools to fit OpenAI tools.

@return A list of OpenAI tools
*/
func convertToOpenAITools() []openai.Tool {
	var openAITools []openai.Tool

	for _, t := range Tools {
		openAITools = append(openAITools, openai.Tool{
			Type: "function",
			Function: &openai.FunctionDefinition{
				Name:        t.Name,
				Description: t.Desc,
				Parameters:  t.Params,
			},
		})
	}

	return openAITools
}

func convertToOpenAIToolCalls(t []ToolCall) []openai.ToolCall {
	var array []openai.ToolCall
	for _, call := range t {
		array = append(array, openai.ToolCall{
			Index: call.Index,
			ID:    call.ID,
			Type:  openai.ToolType(call.Type),
			Function: openai.FunctionCall{
				Arguments: call.Function.Arguments,
				Name:      call.Function.Name,
			},
		})
	}
	return array
}

func convertToHistoryToolCalls(t []openai.ToolCall) []ToolCall {
	var array []ToolCall
	for _, call := range t {
		array = append(array, ToolCall{
			Index: call.Index,
			ID:    call.ID,
			Type:  string(call.Type),
			Function: openai.FunctionCall{
				Arguments: call.Function.Arguments,
				Name:      call.Function.Name,
			},
		})
	}
	return array
}

func createOpenAIMessages(con *Conversation) []openai.ChatCompletionMessage {
	history, err := memory(con)
	if err != nil {
		history = []HistoryStruct{{Role: "system", Content: sp}}
	}

	var array []openai.ChatCompletionMessage

	for _, h := range history {
		array = append(array, openai.ChatCompletionMessage{
			Role:       h.Role,
			Content:    h.Content,
			ToolCallID: h.ToolCallID,
			ToolCalls:  convertToOpenAIToolCalls(h.ToolCalls),
		})
	}

	if rag {
		array = append(array, openai.ChatCompletionMessage{
			Role:    "system",
			Content: getRAG(con.History[len(con.History)-1].Content),
		})
	}

	return array
}
