package dto

import "time"

type SignUpRequest struct {
	Login    string `json:"login"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	LoginOrEmail string `json:"login"`
	Password     string `json:"password"`
}

type UserResponse struct {
	ID        int       `json:"id"`
	Login     string    `json:"login"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
