package dto

type OptionRequest struct {
	Text      string `json:"text"`
	IsCorrect bool   `json:"is_correct"`
}

type OptionResponse struct {
	ID         int    `json:"id"`
	QuestionID int    `json:"question_id"`
	Text       string `json:"text"`
}
