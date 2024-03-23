package services

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"

	openai "github.com/sashabaranov/go-openai"
)

var client *openai.Client

func StartOpenAIClient() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apiToken := os.Getenv("OPEN_API_KEY")

	client = openai.NewClient(apiToken)
}

func GetOpenAPIClient() *openai.Client {
	return client
}

func GetGPTResponse(reqMesages []openai.ChatCompletionMessage) (openai.ChatCompletionMessage, error) {
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo,
			Messages: reqMesages,
		},
	)

	if err != nil {
		return openai.ChatCompletionMessage{}, err
	}

	return resp.Choices[0].Message, nil
}
