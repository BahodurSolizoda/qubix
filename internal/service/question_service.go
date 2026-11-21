package service

import (
	"context"
	"errors"
	"strings"

	"qubix/internal/model"
)

type QuestionService struct {
	repo QuestionRepository
}

func NewQuestionService(repo QuestionRepository) *QuestionService {
	return &QuestionService{
		repo: repo,
	}
}

func validateQuestionType(questionType string) error {
	switch strings.TrimSpace(questionType) {
	case "single", "multiple", "text":
		return nil
	default:
		return errors.New("invalid question type")
	}
}

func (s *QuestionService) Create(ctx context.Context, question *model.Question) error {
	if question.QuizID <= 0 {
		return errors.New("invalid quiz ID")
	}

	if strings.TrimSpace(question.Text) == "" {
		return errors.New("question text is required")
	}

	if question.Points <= 0 {
		return errors.New("question points must be greater than zero")
	}

	if err := validateQuestionType(question.QuestionType); err != nil {
		return err
	}

	return s.repo.Create(ctx, question)
}

func (s *QuestionService) GetByID(ctx context.Context, id int) (*model.Question, error) {
	if id <= 0 {
		return nil, errors.New("invalid question ID")
	}

	return s.repo.GetByID(ctx, id)
}

func (s *QuestionService) GetByQuizID(ctx context.Context, quizID int) ([]model.Question, error) {
	if quizID <= 0 {
		return nil, errors.New("invalid quiz ID")
	}

	return s.repo.GetByQuizID(ctx, quizID)
}

func (s *QuestionService) Update(ctx context.Context, question *model.Question) error {
	if question.ID <= 0 {
		return errors.New("invalid question ID")
	}

	if question.QuizID <= 0 {
		return errors.New("invalid quiz ID")
	}

	if strings.TrimSpace(question.Text) == "" {
		return errors.New("question text is required")
	}

	if question.Points <= 0 {
		return errors.New("question points must be greater than zero")
	}

	if err := validateQuestionType(question.QuestionType); err != nil {
		return err
	}

	return s.repo.Update(ctx, question)
}

func (s *QuestionService) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("invalid question ID")
	}

	return s.repo.Delete(ctx, id)
}
