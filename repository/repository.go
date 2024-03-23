package repository

import (
	"fmt"
	"thinkmate/database"
	"thinkmate/model"
)

func GetMessagesByConversationID(m *[]model.Message, conversationID uint) (err error) {
	db := database.GetDB()
	if db == nil {
		fmt.Println("Error: Database connection is nil.")
		return db.Error
	}

	if err := db.Where("conversation_id = ?", conversationID).Find(&m).Error; err != nil {
		return err
	}
	return nil
}

func SaveMessage(m *model.Message) (err error) {
	db := database.GetDB()
	if db == nil {
		fmt.Println("Error: Database connection is nil.")
		return db.Error
	}

	err = db.Create(&m).Error
	if err != nil {
		return err
	}
	return nil
}
