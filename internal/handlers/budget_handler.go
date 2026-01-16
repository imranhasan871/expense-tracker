package handlers

import (
	"database/sql"
	"encoding/json"
	"expense-tracker/internal/models"
	"expense-tracker/internal/repository"
	"net/http"
	"strconv"
	"strings"
)

type BudgetHandler struct {
	repo        repository.BudgetRepository
	expenseRepo repository.ExpenseRepository
}

func NewBudgetHandler(repo repository.BudgetRepository, expenseRepo repository.ExpenseRepository) *BudgetHandler {
	return &BudgetHandler{repo: repo, expenseRepo: expenseRepo}
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

	budgets, err := h.repo.GetAll(year)
	if err != nil {
		h.sendErrorResponse(w, "Database error", err.Error(), http.StatusInternalServerError)
		return
	}

	summary, err := h.repo.GetDashboardSummary(year)
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

	budget, err := h.repo.GetByCategory(categoryID, year)
	if err != nil {
		if err == sql.ErrNoRows {
			h.sendSuccessResponse(w, map[string]interface{}{
				"allocated": 0,
				"spent":     0,
				"remaining": 0,
				"percent":   0,
			}, "No budget found", http.StatusOK)
			return
		}
		h.sendErrorResponse(w, "Database error", err.Error(), http.StatusInternalServerError)
		return
	}

	spent, err := h.expenseRepo.GetYearlyTotal(categoryID, year)
	if err != nil {
		h.sendErrorResponse(w, "Database error", err.Error(), http.StatusInternalServerError)
		return
	}

	remaining := budget.Amount - spent
	var percent float64
	if budget.Amount > 0 {
		percent = (spent / budget.Amount) * 100
	}

	response := map[string]interface{}{
		"allocated": budget.Amount,
		"spent":     spent,
		"remaining": remaining,
		"percent":   percent,
		"is_locked": budget.IsLocked,
	}

	h.sendSuccessResponse(w, response, "", http.StatusOK)
}

func (h *BudgetHandler) HandleMonitoring(w http.ResponseWriter, r *http.Request) {
	yearStr := r.URL.Query().Get("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		year = 2026
	}

	stats, err := h.repo.GetMonitoringData(year)
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

	if err := h.repo.ToggleLock(id, req.IsLocked); err != nil {
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

	budget, err := h.repo.CreateOrUpdate(req.CategoryID, req.Amount, req.Year)
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
