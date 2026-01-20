package handlers

import (
	"encoding/json"
	"expense-tracker/internal/models"
	"expense-tracker/internal/service"
	"net/http"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Name      string `json:"username"`
		DisplayID string `json:"user_display_id"`
		Email     string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendErrorResponse(w, "Invalid request", err.Error(), http.StatusBadRequest)
		return
	}

	// Default role is Executive - Admin can change this later if needed
	user, err := h.userService.CreateUser(req.Name, req.DisplayID, req.Email, models.RoleExecutive)
	if err != nil {
		h.sendErrorResponse(w, "Failed to create user", err.Error(), http.StatusInternalServerError)
		return
	}

	h.sendSuccessResponse(w, user, "User created successfully. Activation email sent.", http.StatusCreated)
}

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	users, err := h.userService.GetAllUsers()
	if err != nil {
		h.sendErrorResponse(w, "Failed to fetch users", err.Error(), http.StatusInternalServerError)
		return
	}

	h.sendSuccessResponse(w, users, "", http.StatusOK)
}

func (h *UserHandler) UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		UserID int             `json:"user_id"`
		Role   models.UserRole `json:"role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendErrorResponse(w, "Invalid request", err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.userService.UpdateUserRole(req.UserID, req.Role); err != nil {
		h.sendErrorResponse(w, "Failed to update role", err.Error(), http.StatusInternalServerError)
		return
	}

	h.sendSuccessResponse(w, nil, "User role updated successfully", http.StatusOK)
}

func (h *UserHandler) sendErrorResponse(w http.ResponseWriter, error string, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Error: error, Message: message})
}

func (h *UserHandler) sendSuccessResponse(w http.ResponseWriter, data interface{}, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(SuccessResponse{Success: true, Data: data, Message: message})
}
