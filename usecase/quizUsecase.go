package usecase

import "thinkmate/model"

type quizUsecase struct {
	quizRepository model.QuizRepository
}

func NewQuizUsecase(quizRepository model.QuizRepository) model.QuizUsecase {
	return &quizUsecase{
		quizRepository: quizRepository,
	}
}

func (qu *quizUsecase) Create(q *model.Quiz) error {
	return qu.quizRepository.Create(q)
}

func (qu *quizUsecase) FetchQuizByID(q *model.Quiz, quizId string) error {
	return qu.quizRepository.FetchQuizByID(q, quizId)
}

func (qu *quizUsecase) FetchQuizByPin(q *model.Quiz, pin string) error {
	return qu.quizRepository.FetchQuizByPin(q, pin)
}

func (qu *quizUsecase) FetchQuizByTeacherId(q *[]model.Quiz, teacherId uint) error {
	return qu.quizRepository.FetchQuizByTeacherId(q, teacherId)
}
