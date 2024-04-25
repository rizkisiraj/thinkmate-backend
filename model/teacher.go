package model

import (
	"gorm.io/gorm"
)

type Teacher struct {
	gorm.Model
	UUID     string `gorm:"not null" json:uuid`
	FullName string `gorm:"not null" json:"full_name" form:"full_name" valid:"required~Your full name is required"`
	Email    string `gorm:"not null" json:"email" form:"email" valid:"required~Your email is required, email~Invalid email format"`
	Password string `gorm:"not null" json:"password" form:"password" valid:"required~Your password is required, minstringlength(6)~Password has to have a minimum length of 6 characters"`
	Quizzes  []Quiz `json:"quizzes"`
}

type TeacherRepository interface {
	Create(t *Teacher) error
	FetchTeacherByEmail(t *Teacher, email string) error
}

type TeacherUsecase interface {
	Register(t *Teacher) error
	Login(t *Teacher, email string, password string) (string, error)
}
