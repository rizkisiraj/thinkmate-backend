package model

import (
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	ConversationID uint
	Role           string `json:"role"`
	Message        string `json:"message"`
}
