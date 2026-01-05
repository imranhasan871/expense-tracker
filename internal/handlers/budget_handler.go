package handlers

import (
	"encoding/json"
	"expense-tracker/internal/models"
	"expense-tracker/internal/repository"
	"net/http"
	"strconv"
)

// BudgetHandler handles HTTP requests for budgets
type BudgetHandler struct {
	repo *repository.BudgetRepository
}

// NewBudgetHandler creates a new budget handler
func NewBudgetHandler(repo *repository.BudgetRepository) *BudgetHandler {
	return &BudgetHandler{repo: repo}
}

// HandleBudgets handles GET and POST for /api/budgets
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

// GetBudgets retrieves budgets and summary
func (h *BudgetHandler) GetBudgets(w http.ResponseWriter, r *http.Request) {
	yearStr := r.URL.Query().Get("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		year = 2026 // Default year
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

// SetBudget creates or updates a budget
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
