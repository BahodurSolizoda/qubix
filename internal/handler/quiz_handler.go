package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"qubix/internal/dto"
	appmiddleware "qubix/internal/middleware"
	"qubix/internal/model"
	"qubix/internal/service"
)

type QuizHandler struct {
	quizService *service.QuizService
}

func NewQuizHandler(quizService *service.QuizService) *QuizHandler {
	return &QuizHandler{
		quizService: quizService,
	}
}

func (h *QuizHandler) Create(w http.ResponseWriter, r *http.Request) {
	role, ok := r.Context().Value(appmiddleware.UserRoleKey).(string)
	if !ok || role != "author" {
		http.Error(w, "only authors can create quizzes", http.StatusForbidden)
		return
	}

	userID, ok := r.Context().Value(appmiddleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "user ID not found", http.StatusUnauthorized)
		return
	}

	var request dto.QuizRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	quiz := &model.Quiz{
		CreatorID:   userID,
		CategoryID:  request.CategoryID,
		Title:       request.Title,
		Description: request.Description,
		Difficulty:  request.Difficulty,
		TimeLimit:   request.TimeLimit,
	}

	if err := h.quizService.Create(r.Context(), quiz); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := toQuizResponse(quiz)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *QuizHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid quiz ID", http.StatusBadRequest)
		return
	}

	quiz, err := h.quizService.GetByID(r.Context(), id)
	if errors.Is(err, model.ErrNotFound) {
		http.Error(w, "quiz not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "failed to get quiz", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(toQuizResponse(quiz))
}

func (h *QuizHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	quizzes, err := h.quizService.GetAll(r.Context())
	if err != nil {
		http.Error(w, "failed to get quizzes", http.StatusInternalServerError)
		return
	}

	response := make([]dto.QuizResponse, 0, len(quizzes))

	for _, quiz := range quizzes {
		response = append(response, toQuizResponse(&quiz))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *QuizHandler) Update(w http.ResponseWriter, r *http.Request) {
	role, ok := r.Context().Value(appmiddleware.UserRoleKey).(string)
	if !ok || role != "author" {
		http.Error(w, "only authors can update quizzes", http.StatusForbidden)
		return
	}

	userID, ok := r.Context().Value(appmiddleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "user ID not found", http.StatusUnauthorized)
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid quiz ID", http.StatusBadRequest)
		return
	}

	existingQuiz, err := h.quizService.GetByID(r.Context(), id)
	if errors.Is(err, model.ErrNotFound) {
		http.Error(w, "quiz not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "failed to get quiz", http.StatusInternalServerError)
		return
	}

	if existingQuiz.CreatorID != userID {
		http.Error(w, "you can update only your own quizzes", http.StatusForbidden)
		return
	}

	var request dto.QuizRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	quiz := &model.Quiz{
		ID:          id,
		CreatorID:   userID,
		CategoryID:  request.CategoryID,
		Title:       request.Title,
		Description: request.Description,
		Difficulty:  request.Difficulty,
		TimeLimit:   request.TimeLimit,
	}

	if err := h.quizService.Update(r.Context(), quiz); errors.Is(err, model.ErrNotFound) {
		http.Error(w, "quiz not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *QuizHandler) Delete(w http.ResponseWriter, r *http.Request) {
	role, ok := r.Context().Value(appmiddleware.UserRoleKey).(string)
	if !ok || role != "author" {
		http.Error(w, "only authors can delete quizzes", http.StatusForbidden)
		return
	}

	userID, ok := r.Context().Value(appmiddleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "user ID not found", http.StatusUnauthorized)
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid quiz ID", http.StatusBadRequest)
		return
	}

	quiz, err := h.quizService.GetByID(r.Context(), id)
	if errors.Is(err, model.ErrNotFound) {
		http.Error(w, "quiz not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "failed to get quiz", http.StatusInternalServerError)
		return
	}

	if quiz.CreatorID != userID {
		http.Error(w, "you can delete only your own quizzes", http.StatusForbidden)
		return
	}

	if err := h.quizService.Delete(r.Context(), id); errors.Is(err, model.ErrNotFound) {
		http.Error(w, "quiz not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "failed to delete quiz", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func toQuizResponse(quiz *model.Quiz) dto.QuizResponse {
	return dto.QuizResponse{
		ID:          quiz.ID,
		CreatorID:   quiz.CreatorID,
		CategoryID:  quiz.CategoryID,
		Title:       quiz.Title,
		Description: quiz.Description,
		Difficulty:  quiz.Difficulty,
		TimeLimit:   quiz.TimeLimit,
		CreatedAt:   quiz.CreatedAt,
		UpdatedAt:   quiz.UpdatedAt,
	}
}
