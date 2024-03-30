package router

import (
	"thinkmate/controller"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "DELETE", "POST"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://github.com"
		},
		MaxAge: 12 * time.Hour,
	}))

	v1 := r.Group("/v1")
	{
		v1.POST("quiz", controller.CreatQuiz)
		v1.POST("conversation", controller.StartConversation)
		v1.POST("conversation/:id/message", controller.PostAnswer)
		v1.GET("quiz", controller.GetQuizByPin)
	}

	return r
}
