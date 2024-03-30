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
		AllowOrigins:     []string{"*"},                                                // Allows all origins
		AllowMethods:     []string{"PUT", "PATCH", "GET", "DELETE", "POST", "OPTIONS"}, // Ensure OPTIONS is included for preflight requests
		AllowHeaders:     []string{"*"},                                                // Allows all headers. You might specify actual headers you expect or use "*" for any.
		ExposeHeaders:    []string{"Content-Length"},                                   // Specify which headers are safe to expose to the browser
		AllowCredentials: true,                                                         // If you want to include credentials like cookies, authorization headers, or TLS client certificates
		AllowOriginFunc:  func(origin string) bool { return true },                     // Optionally, remove or set to always return true
		MaxAge:           12 * time.Hour,                                               // Defines the max age for the CORS preflight request cache
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
