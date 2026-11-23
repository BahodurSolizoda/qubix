package service

import (
	"context"
	"errors"

	"qubix/internal/model"
)

type ResultService struct {
	attemptRepo  AttemptRepository
	answerRepo   AnswerRepository
	questionRepo QuestionRepository
	optionRepo   OptionRepository
	resultRepo   ResultRepository
}

func NewResultService(attemptRepo AttemptRepository, answerRepo AnswerRepository, questionRepo QuestionRepository, optionRepo OptionRepository, resultRepo ResultRepository) *ResultService {
	return &ResultService{
		attemptRepo:  attemptRepo,
		answerRepo:   answerRepo,
		questionRepo: questionRepo,
		optionRepo:   optionRepo,
		resultRepo:   resultRepo,
	}
}

func (s *ResultService) Finish(ctx context.Context, attemptID, userID int) (*model.Result, error) {
	if attemptID <= 0 || userID <= 0 {
		return nil, errors.New("invalid ID")
	}

	attempt, err := s.attemptRepo.GetByID(ctx, attemptID)
	if err != nil {
		return nil, err
	}

	if attempt.UserID != userID {
		return nil, model.ErrForbidden
	}

	if attempt.Status != "in_progress" {
		return nil, model.ErrAttemptCompleted
	}

	questions, err := s.questionRepo.GetByQuizID(ctx, attempt.QuizID)
	if err != nil {
		return nil, err
	}

	if len(questions) == 0 {
		return nil, errors.New("quiz has no questions")
	}

	answers, err := s.answerRepo.GetByAttemptID(ctx, attemptID)
	if err != nil {
		return nil, err
	}

	answerMap := make(map[int]int)

	for _, answer := range answers {
		answerMap[answer.QuestionID] = answer.ID
	}

	result := &model.Result{
		AttemptID:      attemptID,
		TotalQuestions: len(questions),
	}

	totalPoints := 0

	for _, question := range questions {
		totalPoints += question.Points

		answerID, exists := answerMap[question.ID]
		if !exists {
			continue
		}

		options, err := s.optionRepo.GetByQuestionID(ctx, question.ID)
		if err != nil {
			return nil, err
		}

		selectedIDs, err := s.answerRepo.GetOptionIDsByAnswerID(ctx, answerID)
		if err != nil {
			return nil, err
		}

		correctIDs := make(map[int]bool)

		for _, option := range options {
			if option.IsCorrect {
				correctIDs[option.ID] = true
			}
		}

		if len(correctIDs) == 0 || len(selectedIDs) != len(correctIDs) {
			continue
		}

		correct := true

		for _, id := range selectedIDs {
			if !correctIDs[id] {
				correct = false
				break
			}
		}

		if correct {
			result.CorrectAnswers++
			result.Score += question.Points
		}
	}

	result.IsPassed = result.Score*2 >= totalPoints

	if err := s.resultRepo.CreateAndCompleteAttempt(ctx, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *ResultService) GetResult(ctx context.Context, attemptID, userID int) (*model.Result, error) {
	if attemptID <= 0 || userID <= 0 {
		return nil, errors.New("invalid ID")
	}

	attempt, err := s.attemptRepo.GetByID(ctx, attemptID)
	if err != nil {
		return nil, err
	}

	if attempt.UserID != userID {
		return nil, model.ErrForbidden
	}

	if attempt.Status != "completed" {
		return nil, errors.New("attempt is not completed")
	}

	return s.resultRepo.GetByAttemptID(ctx, attemptID)
}
