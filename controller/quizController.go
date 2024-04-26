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
	"github.com/golang-jwt/jwt/v5"
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
	repository.SaveMessage(&studentMessage)
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
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
		return
	}

	var quiz model.Quiz
	err = repository.GetQuizById(&quiz, postRequest.QuizID)
	if err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
		return
	}

	promptMessage := model.Message{
		ConversationID: newConversation.ID,
		Role:           "system",
		Message:        fmt.Sprintf(prompt, quiz.Topic),
	}

	messages := []model.Message{
		promptMessage,
	}

	message, err := services.GetGPTResponse(messages)
	if err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusBadGateway, gin.H{
			"message": err,
		})
		return
	}

	message.ConversationID = promptMessage.ConversationID
	messagesToSend := []model.Message{
		promptMessage, message,
	}
	err = repository.SaveMessages(&messagesToSend)
	if err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{
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

func GetAllMessagesByConversationId(ctx *gin.Context) {
	id := ctx.Param("id")

	var allMessages []model.Message

	err := repository.GetMessagesByConversationID(&allMessages, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": allMessages,
	})

}

type QuizController struct {
	QuizUsecase model.QuizUsecase
}

func (qc *QuizController) Create(ctx *gin.Context) {
	userData := ctx.MustGet("userData").(jwt.MapClaims)
	userID := uint(userData["id"].(float64))

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
		Topic:     postRequest.Topic,
		Pin:       randomString,
		TeacherID: userID,
	}

	err := qc.QuizUsecase.Create(&newQuiz)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": newQuiz,
	})
	return

}

func (qc *QuizController) FetchById(ctx *gin.Context) {
	id := ctx.Param("id")

	var quiz model.Quiz

	err := qc.QuizUsecase.FetchQuizByID(&quiz, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	if quiz.Topic == "" {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "No matching records",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": quiz,
	})
	return
}

func (qc *QuizController) FetchByPin(ctx *gin.Context) {
	pin, _ := ctx.GetQuery("pin")

	var quiz model.Quiz
	err := qc.QuizUsecase.FetchQuizByPin(&quiz, pin)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusOK, gin.H{
				"status":  http.StatusNotFound,
				"message": "No matching record.",
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
		"data": quiz,
	})
}
