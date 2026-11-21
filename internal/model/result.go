package model

import "time"

type Result struct {
	ID             int       `json:"id"`
	AttemptID      int       `json:"attempt_id"`
	Score          int       `json:"score"`
	TotalQuestions int       `json:"total_questions"`
	CorrectAnswers int       `json:"correct_answers"`
	IsPassed       bool      `json:"is_passed"`
	CreatedAt      time.Time `json:"created_at"`
}
