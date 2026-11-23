package model

import "time"

type Quiz struct {
	ID          int
	CreatorID   int
	CategoryID  int
	Title       string
	Description string
	Difficulty  string
	TimeLimit   int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
