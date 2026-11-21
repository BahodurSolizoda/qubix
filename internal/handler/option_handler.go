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

type OptionHandler struct {
	optionService   *service.OptionService
	questionService *service.QuestionService
	quizService     *service.QuizService
}

type optionRequest struct {
	QuestionID int    `json:"question_id"`
	Text       string `json:"text"`
	IsCorrect  bool   `json:"is_correct"`
}

func NewOptionHandler(
	optionService *service.OptionService,
	questionService *service.QuestionService,
	quizService *service.QuizService,
) *OptionHandler {
	return &OptionHandler{
		optionService:   optionService,
		questionService: questionService,
		quizService:     quizService,
	}
}

func (h *OptionHandler) Create(w http.ResponseWriter, r *http.Request) {
	questionID, err := strconv.Atoi(r.PathValue("question_id"))
	if err != nil {
		http.Error(w, "invalid question ID", http.StatusBadRequest)
		return
	}

	role, ok := r.Context().Value(appmiddleware.UserRoleKey).(string)
	if !ok || role != "author" {
		http.Error(w, "only authors can create options", http.StatusForbidden)
		return
	}

	userID, ok := r.Context().Value(appmiddleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if !h.checkQuizOwnership(w, r, questionID, userID) {
		return
	}

	var request optionRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	option := model.Option{
		QuestionID: questionID,
		Text:       request.Text,
		IsCorrect:  request.IsCorrect,
	}

	if err := h.optionService.Create(r.Context(), &option); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(option)
}

func (h *OptionHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid option ID", http.StatusBadRequest)
		return
	}

	option, err := h.optionService.GetByID(r.Context(), id)
	if errors.Is(err, model.ErrNotFound) {
		http.Error(w, "option not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "failed to get option", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(option)
}

func (h *OptionHandler) GetByQuestionID(w http.ResponseWriter, r *http.Request) {
	questionID, err := strconv.Atoi(r.PathValue("question_id"))
	if err != nil {
		http.Error(w, "invalid question ID", http.StatusBadRequest)
		return
	}

	options, err := h.optionService.GetByQuestionID(r.Context(), questionID)
	if err != nil {
		http.Error(w, "failed to get options", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(options)
}

func (h *OptionHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid option ID", http.StatusBadRequest)
		return
	}

	role, ok := r.Context().Value(appmiddleware.UserRoleKey).(string)
	if !ok || role != "author" {
		http.Error(w, "only authors can update options", http.StatusForbidden)
		return
	}

	userID, ok := r.Context().Value(appmiddleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	existingOption, err := h.optionService.GetByID(r.Context(), id)
	if errors.Is(err, model.ErrNotFound) {
		http.Error(w, "option not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "failed to get option", http.StatusInternalServerError)
		return
	}

	if !h.checkQuizOwnership(w, r, existingOption.QuestionID, userID) {
		return
	}

	var request optionRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	option := model.Option{
		ID:         id,
		QuestionID: existingOption.QuestionID,
		Text:       request.Text,
		IsCorrect:  request.IsCorrect,
	}

	if err := h.optionService.Update(r.Context(), &option); errors.Is(err, model.ErrNotFound) {
		http.Error(w, "option not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *OptionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid option ID", http.StatusBadRequest)
		return
	}

	role, ok := r.Context().Value(appmiddleware.UserRoleKey).(string)
	if !ok || role != "author" {
		http.Error(w, "only authors can delete options", http.StatusForbidden)
		return
	}

	userID, ok := r.Context().Value(appmiddleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	existingOption, err := h.optionService.GetByID(r.Context(), id)
	if errors.Is(err, model.ErrNotFound) {
		http.Error(w, "option not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "failed to get option", http.StatusInternalServerError)
		return
	}

	if !h.checkQuizOwnership(w, r, existingOption.QuestionID, userID) {
		return
	}

	if err := h.optionService.Delete(r.Context(), id); errors.Is(err, model.ErrNotFound) {
		http.Error(w, "option not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "failed to delete option", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *OptionHandler) checkQuizOwnership(
	w http.ResponseWriter,
	r *http.Request,
	questionID int,
	userID int,
) bool {
	question, err := h.questionService.GetByID(r.Context(), questionID)
	if errors.Is(err, model.ErrNotFound) {
		http.Error(w, "question not found", http.StatusNotFound)
		return false
	}
	if err != nil {
		http.Error(w, "failed to get question", http.StatusInternalServerError)
		return false
	}

	quiz, err := h.quizService.GetByID(r.Context(), question.QuizID)
	if errors.Is(err, model.ErrNotFound) {
		http.Error(w, "quiz not found", http.StatusNotFound)
		return false
	}
	if err != nil {
		http.Error(w, "failed to get quiz", http.StatusInternalServerError)
		return false
	}

	if quiz.CreatorID != userID {
		http.Error(w, "you can only manage options in your own quizzes", http.StatusForbidden)
		return false
	}

	return true
}
