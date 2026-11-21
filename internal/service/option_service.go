package service

import (
	"context"
	"errors"
	"strings"

	"qubix/internal/model"
)

type OptionService struct {
	repo OptionRepository
}

func NewOptionService(repo OptionRepository) *OptionService {
	return &OptionService{
		repo: repo,
	}
}

func (s *OptionService) Create(ctx context.Context, option *model.Option) error {
	if option.QuestionID <= 0 {
		return errors.New("invalid question ID")
	}

	if strings.TrimSpace(option.Text) == "" {
		return errors.New("option text is required")
	}

	return s.repo.Create(ctx, option)
}

func (s *OptionService) GetByID(ctx context.Context, id int) (*model.Option, error) {
	if id <= 0 {
		return nil, errors.New("invalid option ID")
	}

	return s.repo.GetByID(ctx, id)
}

func (s *OptionService) GetByQuestionID(ctx context.Context, questionID int) ([]model.Option, error) {
	if questionID <= 0 {
		return nil, errors.New("invalid question ID")
	}

	return s.repo.GetByQuestionID(ctx, questionID)
}

func (s *OptionService) Update(ctx context.Context, option *model.Option) error {
	if option.ID <= 0 {
		return errors.New("invalid option ID")
	}

	if option.QuestionID <= 0 {
		return errors.New("invalid question ID")
	}

	if strings.TrimSpace(option.Text) == "" {
		return errors.New("option text is required")
	}

	return s.repo.Update(ctx, option)
}

func (s *OptionService) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("invalid option ID")
	}

	return s.repo.Delete(ctx, id)
}
