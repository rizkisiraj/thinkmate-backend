package model

import (
	"gorm.io/gorm"
)

type Quiz struct {
	gorm.Model
	Topic         string         `gorm:"not null" json:"topic"`
	Pin           string         `gorm:"type:varchar(4)" json:"pin"`
	Conversations []Conversation `json:"conversations"`
	TeacherID     uint           `json:"teacher_id"`
}

type QuizRepository interface {
	Create(q *Quiz) error
	FetchQuizByID(q *Quiz, quizId string) error
	FetchQuizByPin(q *Quiz, pin string) error
	FetchQuizByTeacherId(q *[]Quiz, teacherId uint) error
}

type QuizUsecase interface {
	Create(q *Quiz) error
	FetchQuizByID(q *Quiz, quizId string) error
	FetchQuizByPin(q *Quiz, pin string) error
	FetchQuizByTeacherId(q *[]Quiz, teacherId uint) error
}
