package model

import "time"

type Result struct {
	ID             int
	AttemptID      int
	Score          int
	TotalQuestions int
	CorrectAnswers int
	IsPassed       bool
	CreatedAt      time.Time
}
