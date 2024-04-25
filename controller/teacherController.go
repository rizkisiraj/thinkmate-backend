package controller

import (
	"net/http"
	helper "thinkmate/helpers"
	"thinkmate/model"

	"github.com/gin-gonic/gin"
)

var (
	appJSON = "application/json"
)

type TeacherController struct {
	TeacherUsecase model.TeacherUsecase
}

func (tc *TeacherController) UserRegister(ctx *gin.Context) {
	contentType := helper.GetContentType(ctx)
	User := model.Teacher{}

	if contentType == appJSON {
		ctx.ShouldBindJSON(&User)
	} else {
		ctx.ShouldBind(&User)
	}

	err := tc.TeacherUsecase.Register(&User)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    User,
	})
}

func (tc *TeacherController) UserLogin(ctx *gin.Context) {
	contentType := helper.GetContentType(ctx)
	User := model.Teacher{}
	var password string

	if contentType == appJSON {
		ctx.ShouldBindJSON(&User)
	} else {
		ctx.ShouldBind(&User)
	}

	password = User.Password

	token, err := tc.TeacherUsecase.Login(&User, User.Email, password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
