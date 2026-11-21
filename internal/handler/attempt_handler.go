package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"qubix/internal/model"
	"qubix/internal/service"
	appmiddleware "qubix/middleware"
)

type AttemptHandler struct {
	attemptService *service.AttemptService
	answerService  *service.AnswerService
	resultService  *service.ResultService
}

func NewAttemptHandler(
	attemptService *service.AttemptService,
	answerService *service.AnswerService,
	resultService *service.ResultService,
) *AttemptHandler {
	return &AttemptHandler{
		attemptService: attemptService,
		answerService:  answerService,
		resultService:  resultService,
	}
}

func (h *AttemptHandler) Start(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(appmiddleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "user ID not found", http.StatusUnauthorized)
		return
	}

	quizID, err := strconv.Atoi(r.PathValue("quiz_id"))
	if err != nil {
		http.Error(w, "invalid quiz ID", http.StatusBadRequest)
		return
	}

	attempt, err := h.attemptService.Start(r.Context(), userID, quizID)
	if errors.Is(err, model.ErrNotFound) {
		http.Error(w, "quiz not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(attempt)
}

func (h *AttemptHandler) SubmitAnswer(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(appmiddleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "user ID not found", http.StatusUnauthorized)
		return
	}

	attemptID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid attempt ID", http.StatusBadRequest)
		return
	}

	var request struct {
		QuestionID int    `json:"question_id"`
		TextAnswer string `json:"text_answer"`
		OptionIDs  []int  `json:"option_ids"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err = h.answerService.Submit(
		r.Context(),
		attemptID,
		userID,
		request.QuestionID,
		request.TextAnswer,
		request.OptionIDs,
	)

	switch {
	case errors.Is(err, model.ErrNotFound):
		http.Error(w, "attempt or question not found", http.StatusNotFound)
		return
	case errors.Is(err, model.ErrForbidden):
		http.Error(w, "you cannot answer another user's attempt", http.StatusForbidden)
		return
	case errors.Is(err, model.ErrAttemptCompleted):
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	case err != nil:
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *AttemptHandler) Finish(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(appmiddleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "user ID not found", http.StatusUnauthorized)
		return
	}

	attemptID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid attempt ID", http.StatusBadRequest)
		return
	}

	result, err := h.resultService.Finish(r.Context(), attemptID, userID)

	switch {
	case errors.Is(err, model.ErrNotFound):
		http.Error(w, "attempt not found", http.StatusNotFound)
		return
	case errors.Is(err, model.ErrForbidden):
		http.Error(w, "you cannot finish another user's attempt", http.StatusForbidden)
		return
	case errors.Is(err, model.ErrAttemptCompleted):
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	case err != nil:
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *AttemptHandler) GetResult(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(appmiddleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "user ID not found", http.StatusUnauthorized)
		return
	}

	attemptID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid attempt ID", http.StatusBadRequest)
		return
	}

	result, err := h.resultService.GetResult(r.Context(), attemptID, userID)

	switch {
	case errors.Is(err, model.ErrNotFound):
		http.Error(w, "attempt or result not found", http.StatusNotFound)
		return
	case errors.Is(err, model.ErrForbidden):
		http.Error(w, "you cannot view another user's result", http.StatusForbidden)
		return
	case err != nil:
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
