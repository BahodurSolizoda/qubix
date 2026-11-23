package dto

import "time"

type AnswerResponse struct {
	ID         int       `json:"id"`
	AttemptID  int       `json:"attempt_id"`
	QuestionID int       `json:"question_id"`
	TextAnswer string    `json:"text_answer"`
	CreatedAt  time.Time `json:"created_at"`
}
