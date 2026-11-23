package service

import (
	"context"
	"errors"

	"qubix/internal/model"
)

type AnswerService struct {
	attemptRepo  AttemptRepository
	questionRepo QuestionRepository
	optionRepo   OptionRepository
	answerRepo   AnswerRepository
}

func NewAnswerService(attemptRepo AttemptRepository, questionRepo QuestionRepository, optionRepo OptionRepository, answerRepo AnswerRepository) *AnswerService {
	return &AnswerService{
		attemptRepo:  attemptRepo,
		questionRepo: questionRepo,
		optionRepo:   optionRepo,
		answerRepo:   answerRepo,
	}
}

func (s *AnswerService) Submit(ctx context.Context, attemptID, userID, questionID int, textAnswer string, optionIDs []int) error {
	if attemptID <= 0 || userID <= 0 || questionID <= 0 {
		return errors.New("invalid ID")
	}

	attempt, err := s.attemptRepo.GetByID(ctx, attemptID)
	if err != nil {
		return err
	}

	if attempt.UserID != userID {
		return model.ErrForbidden
	}

	if attempt.Status != "in_progress" {
		return model.ErrAttemptCompleted
	}

	question, err := s.questionRepo.GetByID(ctx, questionID)
	if err != nil {
		return err
	}

	if question.QuizID != attempt.QuizID {
		return errors.New("question does not belong to this quiz")
	}

	if len(optionIDs) == 0 && textAnswer == "" {
		return errors.New("answer is required")
	}

	options, err := s.optionRepo.GetByQuestionID(ctx, questionID)
	if err != nil {
		return err
	}

	validOptions := make(map[int]bool)

	for _, option := range options {
		validOptions[option.ID] = true
	}

	seen := make(map[int]bool)

	for _, optionID := range optionIDs {
		if !validOptions[optionID] {
			return errors.New("invalid option ID")
		}

		if seen[optionID] {
			return errors.New("duplicate option ID")
		}

		seen[optionID] = true
	}

	answer := &model.Answer{
		AttemptID:  attemptID,
		QuestionID: questionID,
		TextAnswer: textAnswer,
	}

	return s.answerRepo.CreateWithOptions(ctx, answer, optionIDs)
}
