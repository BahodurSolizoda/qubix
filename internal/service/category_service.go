package service

import (
	"context"
	"errors"
	"strings"

	"qubix/internal/model"
)

type CategoryService struct {
	repo CategoryRepository
}

func NewCategoryService(repo CategoryRepository) *CategoryService {
	return &CategoryService{
		repo: repo,
	}
}

func (s *CategoryService) Create(ctx context.Context, category *model.Category) error {
	if strings.TrimSpace(category.Name) == "" {
		return errors.New("category name is required")
	}

	return s.repo.Create(ctx, category)
}

func (s *CategoryService) GetByID(ctx context.Context, id int) (*model.Category, error) {
	if id <= 0 {
		return nil, errors.New("invalid category ID")
	}

	return s.repo.GetByID(ctx, id)
}

func (s *CategoryService) GetAll(ctx context.Context) ([]model.Category, error) {
	return s.repo.GetAll(ctx)
}

func (s *CategoryService) Update(ctx context.Context, category *model.Category) error {
	if category.ID <= 0 {
		return errors.New("invalid category ID")
	}

	if strings.TrimSpace(category.Name) == "" {
		return errors.New("category name is required")
	}

	return s.repo.Update(ctx, category)
}

func (s *CategoryService) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("invalid category ID")
	}

	return s.repo.Delete(ctx, id)
}
