package model

import "time"

type Quiz struct {
	ID          int       `json:"id"`
	CreatorID   int       `json:"creator_id"`
	CategoryID  int       `json:"category_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Difficulty  string    `json:"difficulty"`
	TimeLimit   int       `json:"time_limit"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
