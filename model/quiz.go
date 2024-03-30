package model

import "gorm.io/gorm"

type Quiz struct {
	gorm.Model
	Topic         string         `gorm:"not null" json:"topic"`
	Pin           string         `gorm:"type:varchar(4)" json:"pin"`
	Conversations []Conversation `json:"conversations"`
}
