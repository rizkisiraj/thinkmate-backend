package router

import (
	"thinkmate/controller"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/v1")
	{
		v1.POST("quiz", controller.CreatQuiz)
		v1.POST("conversation", controller.StartConversation)
		v1.POST("conversation/:id/message", controller.PostAnswer)
		v1.GET("quiz", controller.GetQuizByPin)
	}

	return r
}
