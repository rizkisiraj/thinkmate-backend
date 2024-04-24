package repositories

import (
	"thinkmate/model"

	"gorm.io/gorm"
)

type teacherRepository struct {
	database gorm.DB
}

func NewTeacherRepository(db gorm.DB) model.TeacherRepository {
	return &teacherRepository{
		database: db,
	}
}

func (tr *teacherRepository) Create(t *model.Teacher) error {
	err := tr.database.Debug().Create(&t).Error
	return err
}

func (tr *teacherRepository) FetchTeacherByEmail(t *model.Teacher, email string) error {
	err := tr.database.Debug().Where("email = ?", t.Email).Take(&t).Error
	return err
}
