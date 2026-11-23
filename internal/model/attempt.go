package model

import "time"

type Attempt struct {
	ID          int
	UserID      int
	QuizID      int
	StartedAt   time.Time
	CompletedAt *time.Time
	Status      string
}
