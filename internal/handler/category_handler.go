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

type CategoryHandler struct {
	categoryService *service.CategoryService
}

func NewCategoryHandler(categoryService *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

func (h *CategoryHandler) requireAuthor(w http.ResponseWriter, r *http.Request) bool {
	role, ok := r.Context().Value(appmiddleware.UserRoleKey).(string)
	if !ok || role != "author" {
		http.Error(w, "only authors can manage categories", http.StatusForbidden)
		return false
	}

	return true
}

func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	if !h.requireAuthor(w, r) {
		return
	}

	var request dto.CategoryRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	category := &model.Category{
		Name:        request.Name,
		Description: request.Description,
	}

	if err := h.categoryService.Create(r.Context(), category); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := dto.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *CategoryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid category ID", http.StatusBadRequest)
		return
	}

	category, err := h.categoryService.GetByID(r.Context(), id)
	if errors.Is(err, model.ErrNotFound) {
		http.Error(w, "category not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "failed to get category", http.StatusInternalServerError)
		return
	}

	response := dto.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *CategoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	categories, err := h.categoryService.GetAll(r.Context())
	if err != nil {
		http.Error(w, "failed to get categories", http.StatusInternalServerError)
		return
	}

	response := make([]dto.CategoryResponse, 0, len(categories))

	for _, category := range categories {
		response = append(response, dto.CategoryResponse{
			ID:          category.ID,
			Name:        category.Name,
			Description: category.Description,
			CreatedAt:   category.CreatedAt,
			UpdatedAt:   category.UpdatedAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	if !h.requireAuthor(w, r) {
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid category ID", http.StatusBadRequest)
		return
	}

	var request dto.CategoryRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	category := &model.Category{
		ID:          id,
		Name:        request.Name,
		Description: request.Description,
	}

	if err := h.categoryService.Update(r.Context(), category); errors.Is(err, model.ErrNotFound) {
		http.Error(w, "category not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if !h.requireAuthor(w, r) {
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid category ID", http.StatusBadRequest)
		return
	}

	if err := h.categoryService.Delete(r.Context(), id); errors.Is(err, model.ErrNotFound) {
		http.Error(w, "category not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "failed to delete category", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
