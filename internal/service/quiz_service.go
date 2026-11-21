package service

import (
	"context"
	"errors"
	"strings"

	"qubix/internal/model"
)

type QuizService struct {
	repo QuizRepository
}

func NewQuizService(repo QuizRepository) *QuizService {
	return &QuizService{
		repo: repo,
	}
}

func (s *QuizService) Create(ctx context.Context, quiz *model.Quiz) error {
	if quiz.CreatorID <= 0 {
		return errors.New("invalid creator ID")
	}

	if quiz.CategoryID <= 0 {
		return errors.New("invalid category ID")
	}

	if strings.TrimSpace(quiz.Title) == "" {
		return errors.New("quiz title is required")
	}

	if strings.TrimSpace(quiz.Difficulty) == "" {
		return errors.New("quiz difficulty is required")
	}

	if quiz.Difficulty != "easy" &&
		quiz.Difficulty != "medium" &&
		quiz.Difficulty != "hard" {
		return errors.New("invalid quiz difficulty")
	}

	if quiz.TimeLimit < 0 {
		return errors.New("time limit cannot be negative")
	}

	return s.repo.Create(ctx, quiz)
}

func (s *QuizService) GetByID(ctx context.Context, id int) (*model.Quiz, error) {
	if id <= 0 {
		return nil, errors.New("invalid quiz ID")
	}

	return s.repo.GetByID(ctx, id)
}

func (s *QuizService) GetAll(ctx context.Context) ([]model.Quiz, error) {
	return s.repo.GetAll(ctx)
}

func (s *QuizService) Update(ctx context.Context, quiz *model.Quiz) error {
	if quiz.ID <= 0 {
		return errors.New("invalid quiz ID")
	}

	if quiz.CategoryID <= 0 {
		return errors.New("invalid category ID")
	}

	if strings.TrimSpace(quiz.Title) == "" {
		return errors.New("quiz title is required")
	}

	if quiz.Difficulty != "easy" &&
		quiz.Difficulty != "medium" &&
		quiz.Difficulty != "hard" {
		return errors.New("invalid quiz difficulty")
	}

	if quiz.TimeLimit < 0 {
		return errors.New("time limit cannot be negative")
	}

	return s.repo.Update(ctx, quiz)
}

func (s *QuizService) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("invalid quiz ID")
	}

	return s.repo.Delete(ctx, id)
}
