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

type QuestionHandler struct {
	questionService *service.QuestionService
	quizService     *service.QuizService
}

func NewQuestionHandler(
	questionService *service.QuestionService,
	quizService *service.QuizService,
) *QuestionHandler {
	return &QuestionHandler{
		questionService: questionService,
		quizService:     quizService,
	}
}

func (h *QuestionHandler) Create(w http.ResponseWriter, r *http.Request) {
	role, ok := r.Context().Value(appmiddleware.UserRoleKey).(string)
	if !ok || role != "author" {
		http.Error(w, "only authors can create questions", http.StatusForbidden)
		return
	}

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

	quiz, err := h.quizService.GetByID(r.Context(), quizID)
	if errors.Is(err, model.ErrNotFound) {
		http.Error(w, "quiz not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "failed to get quiz", http.StatusInternalServerError)
		return
	}

	if quiz.CreatorID != userID {
		http.Error(w, "you can only add questions to your own quizzes", http.StatusForbidden)
		return
	}

	var question model.Question

	if err := json.NewDecoder(r.Body).Decode(&question); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	question.QuizID = quizID

	if err := h.questionService.Create(r.Context(), &question); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(question)
}

func (h *QuestionHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid question ID", http.StatusBadRequest)
		return
	}

	question, err := h.questionService.GetByID(r.Context(), id)
	if errors.Is(err, model.ErrNotFound) {
		http.Error(w, "question not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "failed to get question", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(question)
}

func (h *QuestionHandler) GetByQuizID(w http.ResponseWriter, r *http.Request) {
	quizID, err := strconv.Atoi(r.PathValue("quiz_id"))
	if err != nil {
		http.Error(w, "invalid quiz ID", http.StatusBadRequest)
		return
	}

	questions, err := h.questionService.GetByQuizID(r.Context(), quizID)
	if err != nil {
		http.Error(w, "failed to get questions", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(questions)
}

func (h *QuestionHandler) Update(w http.ResponseWriter, r *http.Request) {
	role, ok := r.Context().Value(appmiddleware.UserRoleKey).(string)
	if !ok || role != "author" {
		http.Error(w, "only authors can update questions", http.StatusForbidden)
		return
	}

	userID, ok := r.Context().Value(appmiddleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "user ID not found", http.StatusUnauthorized)
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid question ID", http.StatusBadRequest)
		return
	}

	existingQuestion, err := h.questionService.GetByID(r.Context(), id)
	if errors.Is(err, model.ErrNotFound) {
		http.Error(w, "question not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "failed to get question", http.StatusInternalServerError)
		return
	}

	quiz, err := h.quizService.GetByID(r.Context(), existingQuestion.QuizID)
	if errors.Is(err, model.ErrNotFound) {
		http.Error(w, "quiz not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "failed to get quiz", http.StatusInternalServerError)
		return
	}

	if quiz.CreatorID != userID {
		http.Error(w, "you can only update questions in your own quizzes", http.StatusForbidden)
		return
	}

	var question model.Question

	if err := json.NewDecoder(r.Body).Decode(&question); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	question.ID = id
	question.QuizID = existingQuestion.QuizID

	if err := h.questionService.Update(r.Context(), &question); errors.Is(err, model.ErrNotFound) {
		http.Error(w, "question not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *QuestionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	role, ok := r.Context().Value(appmiddleware.UserRoleKey).(string)
	if !ok || role != "author" {
		http.Error(w, "only authors can delete questions", http.StatusForbidden)
		return
	}

	userID, ok := r.Context().Value(appmiddleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "user ID not found", http.StatusUnauthorized)
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid question ID", http.StatusBadRequest)
		return
	}

	question, err := h.questionService.GetByID(r.Context(), id)
	if errors.Is(err, model.ErrNotFound) {
		http.Error(w, "question not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "failed to get question", http.StatusInternalServerError)
		return
	}

	quiz, err := h.quizService.GetByID(r.Context(), question.QuizID)
	if errors.Is(err, model.ErrNotFound) {
		http.Error(w, "quiz not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "failed to get quiz", http.StatusInternalServerError)
		return
	}

	if quiz.CreatorID != userID {
		http.Error(w, "you can only delete questions in your own quizzes", http.StatusForbidden)
		return
	}

	if err := h.questionService.Delete(r.Context(), id); errors.Is(err, model.ErrNotFound) {
		http.Error(w, "question not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "failed to delete question", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
