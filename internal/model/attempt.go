package model

import "time"

type Attempt struct {
	ID          int        `json:"id"`
	UserID      int        `json:"user_id"`
	QuizID      int        `json:"quiz_id"`
	StartedAt   time.Time  `json:"started_at"`
	CompletedAt *time.Time `json:"completed_at"`
	Status      string     `json:"status"`
}
