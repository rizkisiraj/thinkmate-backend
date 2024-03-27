package model

import "gorm.io/gorm"

type Conversation struct {
	gorm.Model
	StudentName string `json:"student_name"`
	QuizID      uint   `json:"quiz_id"`
	Messages    []Message
}
