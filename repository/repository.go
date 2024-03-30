package repository

import (
	"fmt"
	"thinkmate/database"
	"thinkmate/model"
)

func GetMessagesByConversationID(m *[]model.Message, conversationID string) (err error) {
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

func CreateConversation(m *model.Conversation) (err error) {
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

func CreateQuiz(m *model.Quiz) (err error) {
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

func GetQuizByPin(m *model.Quiz, pin string) (err error) {
	db := database.GetDB()
	if db == nil {
		fmt.Println("Error: Database connection is nil.")
		return db.Error
	}

	if err := db.Where("pin = ?", pin).First(&m).Error; err != nil {
		return err
	}
	return nil
}

func GetConversationByQuizId(m *[]model.Conversation, quizId string) (err error) {
	db := database.GetDB()
	if db == nil {
		fmt.Println("Error: Database connection is nil.")
		return db.Error
	}

	if err := db.Where("quiz_id = ?", quizId).Find(&m).Error; err != nil {
		return err
	}
	return nil
}
