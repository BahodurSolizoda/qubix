package model

import "time"

type Answer struct {
	ID         int
	AttemptID  int
	QuestionID int
	TextAnswer string
	CreatedAt  time.Time
}
