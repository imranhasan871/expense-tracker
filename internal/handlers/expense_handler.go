package handlers

import (
	"encoding/json"
	"expense-tracker/internal/models"
	"expense-tracker/internal/service"
	"net/http"
	"strconv"
	"strings"
)

type ExpenseHandler struct {
	service *service.ExpenseService
}

func NewExpenseHandler(service *service.ExpenseService) *ExpenseHandler {
	return &ExpenseHandler{service: service}
}

func (h *ExpenseHandler) HandleExpenses(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if strings.HasSuffix(r.URL.Path, "/insights") || strings.Contains(r.URL.RawQuery, "insights=true") {
			h.GetInsights(w, r)
			return
		}
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

	if err := h.service.Delete(id); err != nil {
		h.sendErrorResponse(w, "Database error", err.Error(), http.StatusInternalServerError)
		return
	}

	h.sendSuccessResponse(w, nil, "Expense deleted successfully", http.StatusOK)
}

func (h *ExpenseHandler) GetExpenses(w http.ResponseWriter, r *http.Request) {
	user := GetAuthenticatedUser(r)
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

	expenses, err := h.service.GetAll(filter, user)
	if err != nil {
		h.sendErrorResponse(w, "Validation error", err.Error(), http.StatusBadRequest)
		return
	}

	h.sendSuccessResponse(w, expenses, "", http.StatusOK)
}

func (h *ExpenseHandler) GetInsights(w http.ResponseWriter, r *http.Request) {
	user := GetAuthenticatedUser(r)
	query := r.URL.Query()
	catID, _ := strconv.Atoi(query.Get("category_id"))

	filter := models.ExpenseFilter{
		StartDate:  query.Get("start_date"),
		EndDate:    query.Get("end_date"),
		CategoryID: catID,
	}

	insights, err := h.service.GetInsights(filter, user)
	if err != nil {
		h.sendErrorResponse(w, "Error fetching insights", err.Error(), http.StatusInternalServerError)
		return
	}

	h.sendSuccessResponse(w, insights, "", http.StatusOK)
}

func (h *ExpenseHandler) CreateExpense(w http.ResponseWriter, r *http.Request) {
	user := GetAuthenticatedUser(r)
	var req models.ExpenseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendErrorResponse(w, "Invalid JSON", err.Error(), http.StatusBadRequest)
		return
	}

	expense, err := h.service.Create(req, user) // Pass user to service method
	if err != nil {
		if err.Error() == "spending is temporarily locked for this category" {
			h.sendErrorResponse(w, "Circuit Breaker Active", err.Error(), http.StatusForbidden)
		} else {
			h.sendErrorResponse(w, "Validation error", err.Error(), http.StatusBadRequest)
		}
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
