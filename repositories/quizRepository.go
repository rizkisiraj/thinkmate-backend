package repositories

import (
	"thinkmate/model"

	"gorm.io/gorm"
)

type quizRepository struct {
	database gorm.DB
}

func NewQuizRepository(db gorm.DB) model.QuizRepository {
	return &quizRepository{
		database: db,
	}
}

func (qr *quizRepository) Create(q *model.Quiz) error {
	err := qr.database.Create(&q).Error
	if err != nil {
		return err
	}
	return nil
}

func (qr *quizRepository) FetchQuizByID(q *model.Quiz, quizId string) error {
	if err := qr.database.Where("id = ?", quizId).Find(&q).Error; err != nil {
		return err
	}
	return nil
}

func (qr *quizRepository) FetchQuizByPin(q *model.Quiz, pin string) error {
	if err := qr.database.Where("pin = ?", pin).First(&q).Error; err != nil {
		return err
	}
	return nil
}

func (qr *quizRepository) FetchQuizByTeacherId(q *[]model.Quiz, teacherId uint) error {
	if err := qr.database.Where("teacherId = ?", teacherId).Find(&q).Error; err != nil {
		return err
	}
	return nil
}
