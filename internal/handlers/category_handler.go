package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"expense-tracker/internal/models"
	"expense-tracker/internal/service"
)

type CategoryHandler struct {
	service *service.CategoryService
}

func NewCategoryHandler(service *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Message string      `json:"message,omitempty"`
}

func (h *CategoryHandler) HandleCategories(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetAllCategories(w, r)
	case http.MethodPost:
		h.CreateCategory(w, r)
	default:
		h.sendErrorResponse(w, "Method not allowed", "Only GET and POST methods are supported", http.StatusMethodNotAllowed)
	}
}

func (h *CategoryHandler) HandleCategoryByID(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	if path == "" {
		h.sendErrorResponse(w, "Invalid request", "Category ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(path)
	if err != nil {
		h.sendErrorResponse(w, "Invalid ID", "Category ID must be a valid number", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.GetCategoryByID(w, r, id)
	case http.MethodPut:
		h.UpdateCategory(w, r, id)
	case http.MethodPatch:
		h.ToggleCategoryStatus(w, r, id)
	default:
		h.sendErrorResponse(w, "Method not allowed", "Supported: GET, PUT, PATCH", http.StatusMethodNotAllowed)
	}
}

func (h *CategoryHandler) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	activeOnly := r.URL.Query().Get("active_only") == "true"

	categories, err := h.service.GetAll(activeOnly)
	if err != nil {
		h.sendErrorResponse(w, "Database error", err.Error(), http.StatusInternalServerError)
		return
	}

	h.sendSuccessResponse(w, categories, "", http.StatusOK)
}

func (h *CategoryHandler) GetCategoryByID(w http.ResponseWriter, r *http.Request, id int) {
	category, err := h.service.GetByID(id)
	if err != nil {
		if err.Error() == "category not found" {
			h.sendErrorResponse(w, "Not found", "Category not found", http.StatusNotFound)
		} else {
			h.sendErrorResponse(w, "Database error", err.Error(), http.StatusInternalServerError)
		}
		return
	}

	h.sendSuccessResponse(w, category, "", http.StatusOK)
}

func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var req models.CategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendErrorResponse(w, "Invalid JSON", "Request body must be valid JSON", http.StatusBadRequest)
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		h.sendErrorResponse(w, "Validation error", "Category name is required", http.StatusBadRequest)
		return
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	category, err := h.service.Create(req.Name, isActive)
	if err != nil {
		if err.Error() == "category with this name already exists" {
			h.sendErrorResponse(w, "Duplicate category", err.Error(), http.StatusConflict)
		} else {
			h.sendErrorResponse(w, "Database error", err.Error(), http.StatusInternalServerError)
		}
		return
	}

	h.sendSuccessResponse(w, category, "Category created successfully", http.StatusCreated)
}

func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request, id int) {
	var req models.CategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendErrorResponse(w, "Invalid JSON", "Request body must be valid JSON", http.StatusBadRequest)
		return
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	category, err := h.service.Update(id, req.Name, isActive)
	if err != nil {
		if err.Error() == "category not found" {
			h.sendErrorResponse(w, "Not found", "Category not found", http.StatusNotFound)
		} else {
			h.sendErrorResponse(w, "Database error", err.Error(), http.StatusInternalServerError)
		}
		return
	}

	h.sendSuccessResponse(w, category, "Category updated successfully", http.StatusOK)
}

func (h *CategoryHandler) ToggleCategoryStatus(w http.ResponseWriter, r *http.Request, id int) {
	category, err := h.service.ToggleStatus(id)
	if err != nil {
		if err.Error() == "category not found" {
			h.sendErrorResponse(w, "Not found", "Category not found", http.StatusNotFound)
		} else {
			h.sendErrorResponse(w, "Database error", err.Error(), http.StatusInternalServerError)
		}
		return
	}

	h.sendSuccessResponse(w, category, "Status toggled successfully", http.StatusOK)
}

func (h *CategoryHandler) sendErrorResponse(w http.ResponseWriter, error string, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error:   error,
		Message: message,
	})
}

func (h *CategoryHandler) sendSuccessResponse(w http.ResponseWriter, data interface{}, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(SuccessResponse{
		Success: true,
		Data:    data,
		Message: message,
	})
}
