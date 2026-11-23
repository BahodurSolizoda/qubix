package dto

import "time"

type SubmitAnswerRequest struct {
	QuestionID int    `json:"question_id"`
	TextAnswer string `json:"text_answer"`
	OptionIDs  []int  `json:"option_ids"`
}

type AttemptResponse struct {
	ID          int        `json:"id"`
	UserID      int        `json:"user_id"`
	QuizID      int        `json:"quiz_id"`
	StartedAt   time.Time  `json:"started_at"`
	CompletedAt *time.Time `json:"completed_at"`
	Status      string     `json:"status"`
}
