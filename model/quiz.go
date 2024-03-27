package model

import "gorm.io/gorm"

type Quiz struct {
	gorm.Model
	Topic         string `gorm:"not null" json:"topic"`
	Conversations []Conversation
}
