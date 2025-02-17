package provider

import (
	"Coeus/llm/tool"
	"context"
	"encoding/json"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

type OpenAIStruct struct {
	Model   string
	API_KEY string
}

func OpenAI(model, api_key string) error {
	Provider = OpenAIStruct{
		Model:   model,
		API_KEY: api_key,
	}
	return nil
}

func SendOpenAI(request RequestStruct) (ResponseStruct, error) {
	config := Provider.(OpenAIStruct)
	client := openai.NewClient(config.API_KEY)

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
			tool, err := tool.Find(t.Function.Name)
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

			toolResponse, err := tool.Run(args...)
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

	for _, t := range tool.Tools {
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
