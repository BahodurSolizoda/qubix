package dto

import "time"

type CategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CategoryResponse struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
