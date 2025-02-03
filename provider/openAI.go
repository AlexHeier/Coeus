package provider

import (
	"context"
	"fmt"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
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
	client := openai.NewClient(option.WithAPIKey(config.API_KEY))

	chatCompletion, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
		}),
		Model: openai.F(openai.ChatModelGPT4o),
	})
	if err != nil {
		return ResponseStruct{}, err
	}
	fmt.Print(chatCompletion.Choices[0].Message.Content)

	return ResponseStruct{
		Response: chatCompletion.Choices[0].Message.Content,
	}, nil
}
