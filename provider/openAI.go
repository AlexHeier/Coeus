package provider

import (
	"context"

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

func SendOpenAI(prompt string) (ResponseStruct, error) {
	config := Provider.(OpenAIStruct)
	client := openai.NewClient(config.API_KEY)

	resp, err := client.CreateChatCompletion(context.TODO(), openai.ChatCompletionRequest{
		Model: config.Model,
		Messages: []openai.ChatCompletionMessage{
			{Role: "user", Content: prompt},
		},
	})
	if err != nil {
		return ResponseStruct{}, err
	}

	responseContent := resp.Choices[0].Message.Content

	return ResponseStruct{Response: responseContent}, nil
}
