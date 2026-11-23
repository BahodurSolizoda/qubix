package service

import (
	"context"
	"errors"

	"qubix/internal/model"
)

type AttemptService struct {
	attemptRepo AttemptRepository
	quizRepo    QuizRepository
}

func NewAttemptService(attemptRepo AttemptRepository, quizRepo QuizRepository) *AttemptService {
	return &AttemptService{
		attemptRepo: attemptRepo,
		quizRepo:    quizRepo,
	}
}

func (s *AttemptService) Start(ctx context.Context, userID, quizID int) (*model.Attempt, error) {
	if userID <= 0 {
		return nil, errors.New("invalid user ID")
	}

	if quizID <= 0 {
		return nil, errors.New("invalid quiz ID")
	}

	_, err := s.quizRepo.GetByID(ctx, quizID)
	if err != nil {
		return nil, err
	}

	attempt := &model.Attempt{
		UserID: userID,
		QuizID: quizID,
		Status: "in_progress",
	}

	if err := s.attemptRepo.Create(ctx, attempt); err != nil {
		return nil, err
	}

	return attempt, nil
}

func (s *AttemptService) GetByID(ctx context.Context, id int) (*model.Attempt, error) {
	if id <= 0 {
		return nil, errors.New("invalid attempt ID")
	}

	return s.attemptRepo.GetByID(ctx, id)
}
