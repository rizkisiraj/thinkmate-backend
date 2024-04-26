package usecase

import (
	"errors"
	helper "thinkmate/helpers"
	"thinkmate/model"

	"github.com/google/uuid"
)

type teacherUsecase struct {
	teacherRepository model.TeacherRepository
}

func NewTeacherUsecase(teacherRepository model.TeacherRepository) model.TeacherUsecase {
	return &teacherUsecase{
		teacherRepository: teacherRepository,
	}
}

func (tu *teacherUsecase) Register(t *model.Teacher) error {
	newUUID := uuid.New()
	t.UUID = newUUID.String()

	t.Password = helper.HashPass(t.Password)

	return tu.teacherRepository.Create(t)
}

func (tu *teacherUsecase) Login(t *model.Teacher, email string, password string) (string, error) {
	err := tu.teacherRepository.FetchTeacherByEmail(t, email)
	if err != nil {
		return "", errors.New("invalid email")
	}

	// Check password - assuming you have a method to compare hashed passwords
	if !helper.ComparePass([]byte(t.Password), []byte(password)) {
		return "", errors.New("invalid password")
	}

	// Generate token
	token := helper.GenerateToken(t.ID, t.Email)
	return token, nil
}
