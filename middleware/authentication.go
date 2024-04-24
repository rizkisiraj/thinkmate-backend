package middleware

import (
	"net/http"
	helper "thinkmate/helpers"

	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		verifyToken, err := helper.VerifyToken(ctx)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unathenticated",
				"message": err.Error(),
			})
			return
		}

		ctx.Set("userData", verifyToken)
		ctx.Next()
	}
}
