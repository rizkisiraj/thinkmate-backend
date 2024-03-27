package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"thinkmate/model"
	"thinkmate/repository"
	"thinkmate/services"

	"github.com/gin-gonic/gin"
)

const prompt = "Kamu adalah ThinkMate AI, Teman diskusi siswa SMA. Kamu dapat memantik diskusi dari topik yang sudah ditentukan guru. Topiknya adalah %s Kamu dapat me-encourage siswa untuk berargumen. Kamu bisa memberi pertanyaan lanjutan dari argumen yang sebelumnya diberi siswa. Kamu dapat memvalidasi benar atau salah pernyataan argument siswa dengan menyocokan fakta pengetahuan dari sumber yang reliable. Jika siswa to the poin bertanya apa jawaban dari suatu hal, jangan langsung diberi jawaban, tapi encourage siswa untuk berpikir apa jawabannya, kamu bisa berikan Langkah Langkah berpikirnya. Berikan HANYA 1 PERTANYAAN DAN JANGAN MEMBUAT PANJANG PERCAKAPAN. JAWABLAH DENGAN SIMPLE."

func PostAnswer(ctx *gin.Context) {
	id := ctx.Param("id")
	var messages []model.Message

	postRequest := struct {
		ConversationID uint   `json:"conversation_id"`
		Message        string `json:"message"`
	}{}

	if err := ctx.BindJSON(&postRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err,
		})
		return
	}

	repository.GetMessagesByConversationID(&messages, id)

	var studentMessage = model.Message{
		ConversationID: postRequest.ConversationID,
		Role:           "Student",
		Message:        postRequest.Message,
	}

	messages = append(messages, studentMessage)

	gptMessage, err := services.GetGPTResponse(messages)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
		return
	}

	conversationId, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
		return
	}

	gptMessage.ConversationID = uint(conversationId)
	repository.SaveMessage(&gptMessage)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Successfully created",
		"data":    gptMessage,
	})
}

func StartConversation(ctx *gin.Context) {
	postRequest := struct {
		QuizID uint   `json:"quiz_id"`
		Name   string `json:"name"`
	}{}

	if err := ctx.BindJSON(&postRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err,
		})
		return
	}

	newConversation := model.Conversation{
		StudentName: postRequest.Name,
		QuizID:      postRequest.QuizID,
	}

	err := repository.CreateConversation(&newConversation)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
		return
	}

	promptMessage := model.Message{
		ConversationID: newConversation.ID,
		Role:           "system",
		Message:        fmt.Sprintf(prompt, "Partisipasi perempuan di bidang IT"),
	}

	messages := []model.Message{
		promptMessage,
	}

	repository.SaveMessage(&promptMessage)

	message, err := services.GetGPTResponse(messages)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"message": err,
		})
		return
	}

	postResponse := struct {
		ConversationID uint   `json:"conversation_id"`
		Message        string `json:"message"`
	}{
		ConversationID: newConversation.ID,
		Message:        message.Message,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Successfully created",
		"data":    postResponse,
	})
}

func CreatQuiz(ctx *gin.Context) {
	postRequest := struct {
		Topic string `json:"topic"`
	}{}

	if err := ctx.BindJSON(&postRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err,
		})
		return
	}

	newQuiz := model.Quiz{
		Topic: postRequest.Topic,
	}

	err := repository.CreateQuiz(&newQuiz)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Successfully created",
		"data":    newQuiz,
	})
}
