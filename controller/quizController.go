package controller

import (
	"net/http"
	"thinkmate/services"

	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
)

func postAnswer(ctx *gin.Context) {
	var messages []openai.ChatCompletionMessage

	if err := ctx.BindJSON(&messages); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err,
		})
		return
	}

	message, err := services.GetGPTResponse(messages)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
		return
	}

	messages = append(messages, message)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Successfully created",
		"data":    messages,
	})

}
