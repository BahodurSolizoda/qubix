package model

import "time"

type Question struct {
	ID           int
	QuizID       int
	Text         string
	QuestionType string
	Points       int
	CreatedAt    time.Time
}
