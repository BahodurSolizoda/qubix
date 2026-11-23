package dto

import "time"

type QuestionRequest struct {
	Text         string `json:"text"`
	QuestionType string `json:"question_type"`
	Points       int    `json:"points"`
}

type QuestionResponse struct {
	ID           int       `json:"id"`
	QuizID       int       `json:"quiz_id"`
	Text         string    `json:"text"`
	QuestionType string    `json:"question_type"`
	Points       int       `json:"points"`
	CreatedAt    time.Time `json:"created_at"`
}
