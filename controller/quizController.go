package controller

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"thinkmate/database"
	"thinkmate/model"
	"thinkmate/repository"
	"thinkmate/services"

	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
	"gorm.io/gorm"
)

const prompt = "Kamu adalah ThinkMate AI, Teman diskusi siswa SMA. Kamu dapat memantik diskusi dari topik yang sudah ditentukan guru. Topiknya adalah %s Kamu dapat me-encourage siswa untuk berargumen. Kamu bisa memberi pertanyaan lanjutan dari argumen yang sebelumnya diberi siswa. Kamu dapat memvalidasi benar atau salah pernyataan argument siswa dengan menyocokan fakta pengetahuan dari sumber yang reliable. Jika siswa to the poin bertanya apa jawaban dari suatu hal, jangan langsung diberi jawaban, tapi encourage siswa untuk berpikir apa jawabannya, kamu bisa berikan Langkah Langkah berpikirnya. Berikan HANYA 1 PERTANYAAN DAN JANGAN MEMBUAT PANJANG PERCAKAPAN. JAWABLAH DENGAN SIMPLE."

func PostAnswer(ctx *gin.Context) {
	id := ctx.Param("id")
	var messages []model.Message

	conversationId, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
		return
	}

	postRequest := struct {
		Message string `json:"message"`
	}{}

	if err := ctx.BindJSON(&postRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err,
		})
		return
	}

	repository.GetMessagesByConversationID(&messages, id)

	var studentMessage = model.Message{
		ConversationID: uint(conversationId),
		Role:           openai.ChatMessageRoleUser,
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

	tx := database.GetDB().Begin()
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

	err = repository.SaveMessage(&promptMessage)
	if err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
		return
	}

	message, err := services.GetGPTResponse(messages)
	if err != nil {
		tx.Rollback()
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

	tx.Commit()

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Successfully created",
		"data":    postResponse,
	})
}

func CreatQuiz(ctx *gin.Context) {
	postRequest := struct {
		Topic string `json:"topic"`
	}{}

	randomNumber := rand.Intn(9000) + 1000
	randomString := fmt.Sprintf("%d", randomNumber)

	if err := ctx.BindJSON(&postRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err,
		})
		return
	}

	newQuiz := model.Quiz{
		Topic: postRequest.Topic,
		Pin:   randomString,
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

func GetQuizByPin(ctx *gin.Context) {
	pin, _ := ctx.GetQuery("pin")

	var quizToSend model.Quiz
	err := repository.GetQuizByPin(&quizToSend, pin)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusOK, gin.H{
				"status":  http.StatusNotFound,
				"message": "No quiz with matching pin.",
			})
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": err,
			})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": quizToSend,
	})
}

func GetAllConversationByQuizId(ctx *gin.Context) {
	id := ctx.Param("id")

	var allConversations []model.Conversation

	err := repository.GetConversationByQuizId(&allConversations, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": allConversations,
	})

}
