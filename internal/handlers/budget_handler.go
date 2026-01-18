package handlers

import (
	"encoding/json"
	"expense-tracker/internal/models"
	"expense-tracker/internal/service"
	"net/http"
	"strconv"
	"strings"
)

type BudgetHandler struct {
	service *service.BudgetService
}

func NewBudgetHandler(service *service.BudgetService) *BudgetHandler {
	return &BudgetHandler{service: service}
}

func (h *BudgetHandler) HandleBudgets(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetBudgets(w, r)
	case http.MethodPost:
		h.SetBudget(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *BudgetHandler) GetBudgets(w http.ResponseWriter, r *http.Request) {
	yearStr := r.URL.Query().Get("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		year = 2026
	}

	budgets, err := h.service.GetAll(year)
	if err != nil {
		h.sendErrorResponse(w, "Database error", err.Error(), http.StatusInternalServerError)
		return
	}

	summary, err := h.service.GetDashboardSummary(year)
	if err != nil {
		h.sendErrorResponse(w, "Database error", err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"budgets": budgets,
		"summary": summary,
	}

	h.sendSuccessResponse(w, data, "", http.StatusOK)
}

func (h *BudgetHandler) GetBudgetStatus(w http.ResponseWriter, r *http.Request) {
	categoryIDStr := r.URL.Query().Get("category_id")
	yearStr := r.URL.Query().Get("year")

	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		h.sendErrorResponse(w, "Invalid category_id", "Category ID is required", http.StatusBadRequest)
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		h.sendErrorResponse(w, "Invalid year", "Year is required", http.StatusBadRequest)
		return
	}

	status, err := h.service.GetStatus(categoryID, year)
	if err != nil {
		h.sendErrorResponse(w, "Database error", err.Error(), http.StatusInternalServerError)
		return
	}

	h.sendSuccessResponse(w, status, "", http.StatusOK)
}

func (h *BudgetHandler) HandleMonitoring(w http.ResponseWriter, r *http.Request) {
	yearStr := r.URL.Query().Get("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		year = 2026
	}

	stats, err := h.service.GetMonitoringData(year)
	if err != nil {
		h.sendErrorResponse(w, "Database error", err.Error(), http.StatusInternalServerError)
		return
	}

	h.sendSuccessResponse(w, stats, "", http.StatusOK)
}

func (h *BudgetHandler) ToggleCircuitBreaker(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		h.sendErrorResponse(w, "Invalid URL", "Budget ID missing", http.StatusBadRequest)
		return
	}
	idStr := pathParts[3]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.sendErrorResponse(w, "Invalid ID", "Budget ID must be a number", http.StatusBadRequest)
		return
	}

	var req struct {
		IsLocked bool `json:"is_locked"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendErrorResponse(w, "Invalid JSON", err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.ToggleLock(id, req.IsLocked); err != nil {
		h.sendErrorResponse(w, "Database error", err.Error(), http.StatusInternalServerError)
		return
	}

	status := "unlocked"
	if req.IsLocked {
		status = "locked"
	}
	h.sendSuccessResponse(w, nil, "Circuit breaker "+status, http.StatusOK)
}

func (h *BudgetHandler) SetBudget(w http.ResponseWriter, r *http.Request) {
	var req models.BudgetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendErrorResponse(w, "Invalid JSON", err.Error(), http.StatusBadRequest)
		return
	}

	if req.CategoryID <= 0 || req.Amount <= 0 || req.Year <= 0 {
		h.sendErrorResponse(w, "Validation error", "CategoryID, Amount, and Year are required", http.StatusBadRequest)
		return
	}

	budget, err := h.service.CreateOrUpdate(req.CategoryID, req.Amount, req.Year)
	if err != nil {
		h.sendErrorResponse(w, "Database error", err.Error(), http.StatusInternalServerError)
		return
	}

	h.sendSuccessResponse(w, budget, "Budget set successfully", http.StatusOK)
}

func (h *BudgetHandler) sendErrorResponse(w http.ResponseWriter, error string, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Error: error, Message: message})
}

func (h *BudgetHandler) sendSuccessResponse(w http.ResponseWriter, data interface{}, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(SuccessResponse{Success: true, Data: data, Message: message})
}
