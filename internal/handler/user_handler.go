package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"qubix/internal/dto"
	"qubix/internal/model"
	"qubix/internal/service"
)

type UserHandler struct {
	service    *service.UserService
	jwtService *service.JWTService
}

func NewUserHandler(userService *service.UserService, jwtService *service.JWTService) *UserHandler {
	return &UserHandler{
		service:    userService,
		jwtService: jwtService,
	}
}

func (h *UserHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var request dto.SignUpRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	user := &model.User{
		Login: request.Login,
		Email: request.Email,
	}

	err := h.service.Register(r.Context(), user, request.Password)
	if errors.Is(err, model.ErrDuplicate) {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := dto.UserResponse{
		ID:        user.ID,
		Login:     user.Login,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var request dto.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.service.Login(r.Context(), request.LoginOrEmail, request.Password)
	if err != nil {
		http.Error(w, model.ErrInvalidCredentials.Error(), http.StatusUnauthorized)
		return
	}

	token, err := h.jwtService.GenerateToken(user.ID, user.Role)
	if err != nil {
		http.Error(w, "could not generate token", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"user": dto.UserResponse{
			ID:        user.ID,
			Login:     user.Login,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
		"token": token,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
