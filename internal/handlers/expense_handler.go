package handlers

import (
	"encoding/json"
	"expense-tracker/internal/models"
	"expense-tracker/internal/repository"
	"net/http"
	"strconv"
	"strings"
)

type ExpenseHandler struct {
	repo       repository.ExpenseRepository
	budgetRepo repository.BudgetRepository
}

func NewExpenseHandler(repo repository.ExpenseRepository, budgetRepo repository.BudgetRepository) *ExpenseHandler {
	return &ExpenseHandler{repo: repo, budgetRepo: budgetRepo}
}

func (h *ExpenseHandler) HandleExpenses(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetExpenses(w, r)
	case http.MethodPost:
		h.CreateExpense(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *ExpenseHandler) HandleExpenseByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/expenses/")
	id, err := strconv.Atoi(path)
	if err != nil {
		h.sendErrorResponse(w, "Invalid ID", "Expense ID must be a number", http.StatusBadRequest)
		return
	}

	if err := h.repo.Delete(id); err != nil {
		h.sendErrorResponse(w, "Database error", err.Error(), http.StatusInternalServerError)
		return
	}

	h.sendSuccessResponse(w, nil, "Expense deleted successfully", http.StatusOK)
}

func (h *ExpenseHandler) GetExpenses(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	catID, _ := strconv.Atoi(query.Get("category_id"))
	minAmount, _ := strconv.ParseFloat(query.Get("min_amount"), 64)
	maxAmount, _ := strconv.ParseFloat(query.Get("max_amount"), 64)

	filter := models.ExpenseFilter{
		StartDate:  query.Get("start_date"),
		EndDate:    query.Get("end_date"),
		CategoryID: catID,
		SearchText: query.Get("search"),
		MinAmount:  minAmount,
		MaxAmount:  maxAmount,
	}

	if err := filter.Validate(); err != nil {
		h.sendErrorResponse(w, "Validation error", err.Error(), http.StatusBadRequest)
		return
	}

	expenses, err := h.repo.GetAll(filter)
	if err != nil {
		h.sendErrorResponse(w, "Database error", err.Error(), http.StatusInternalServerError)
		return
	}

	h.sendSuccessResponse(w, expenses, "", http.StatusOK)
}

func (h *ExpenseHandler) CreateExpense(w http.ResponseWriter, r *http.Request) {
	var req models.ExpenseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendErrorResponse(w, "Invalid JSON", err.Error(), http.StatusBadRequest)
		return
	}

	if req.CategoryID <= 0 || req.Amount <= 0 || req.ExpenseDate == "" {
		h.sendErrorResponse(w, "Validation error", "Category, Amount, and Date are required", http.StatusBadRequest)
		return
	}

	if len(req.ExpenseDate) >= 4 {
		year, err := strconv.Atoi(req.ExpenseDate[:4])
		if err == nil {
			isLocked, err := h.budgetRepo.IsLocked(req.CategoryID, year)
			if err != nil {
				h.sendErrorResponse(w, "Database error", "Failed to check budget lock status", http.StatusInternalServerError)
				return
			}
			if isLocked {
				h.sendErrorResponse(w, "Circuit Breaker Active", "Spending is temporarily locked for this category", http.StatusForbidden)
				return
			}
		}
	}

	expense, err := h.repo.Create(req)
	if err != nil {
		h.sendErrorResponse(w, "Database error", err.Error(), http.StatusInternalServerError)
		return
	}

	h.sendSuccessResponse(w, expense, "Expense recorded successfully", http.StatusCreated)
}

func (h *ExpenseHandler) sendErrorResponse(w http.ResponseWriter, error string, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Error: error, Message: message})
}

func (h *ExpenseHandler) sendSuccessResponse(w http.ResponseWriter, data interface{}, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(SuccessResponse{Success: true, Data: data, Message: message})
}
