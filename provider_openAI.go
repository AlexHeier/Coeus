package coeus

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

func OpenAI(model, apikey string) error {
	Provider = openAIStruct{
		Model:  model,
		ApiKey: apikey,
	}
	return nil
}

func sendOpenAI(request RequestStruct) (ResponseStruct, error) {
	config := Provider.(openAIStruct)
	client := openai.NewClient(config.ApiKey)

	openAITools := convertToOpenAITools()

	resp, err := client.CreateChatCompletion(context.TODO(), openai.ChatCompletionRequest{
		Model: config.Model,
		Messages: []openai.ChatCompletionMessage{
			{Role: "user", Content: request.Systemprompt + request.Userprompt},
		},
		Tools: openAITools,
	})
	if err != nil {
		return ResponseStruct{}, err
	}

	if len(resp.Choices[0].Message.ToolCalls) > 0 {
		var newMessage []openai.ChatCompletionMessage
		newMessage = append(newMessage, resp.Choices[0].Message)
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

			newMessage = append(newMessage, openai.ChatCompletionMessage{
				Role:       "tool",
				ToolCallID: t.ID,
				Content:    string(toolResponse),
			})
		}
		resp, err = client.CreateChatCompletion(context.TODO(), openai.ChatCompletionRequest{
			Model:    config.Model,
			Messages: newMessage,
			Tools:    openAITools,
		})
		if err != nil {
			return ResponseStruct{}, err
		}
	}

	return ResponseStruct{Response: resp.Choices[0].Message.Content}, nil
}

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
