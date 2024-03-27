package services

import (
	"context"
	"log"
	"os"
	helper "thinkmate/helpers"
	"thinkmate/model"

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

func GetGPTResponse(reqMesages []model.Message) (model.Message, error) {
	convertedMessages := helper.ConvertModels(reqMesages)

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:       openai.GPT3Dot5Turbo,
			MaxTokens:   60,
			Temperature: 0.6,
			Messages:    convertedMessages,
		},
	)

	if err != nil {
		return model.Message{}, err
	}

	responseMessage := model.Message{
		Role:    resp.Choices[0].Message.Role,
		Message: resp.Choices[0].Message.Content,
	}

	return responseMessage, nil
}
