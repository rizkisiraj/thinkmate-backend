package helper

import (
	"thinkmate/model"

	"github.com/sashabaranov/go-openai"
)

func ConvertModels(originals []model.Message) []openai.ChatCompletionMessage {
	reduced := make([]openai.ChatCompletionMessage, len(originals))
	for i, o := range originals {
		reduced[i] = openai.ChatCompletionMessage{
			Role:    o.Role,
			Content: o.Message,
		}
	}
	return reduced
}
