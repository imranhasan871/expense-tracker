package handlers

import (
	"encoding/json"
	"expense-tracker/internal/service"
	"net/http"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendErrorResponse(w, "Invalid request", err.Error(), http.StatusBadRequest)
		return
	}

	user, token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		h.sendErrorResponse(w, "Authentication failed", err.Error(), http.StatusUnauthorized)
		return
	}

	// Set cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
	})

	h.sendSuccessResponse(w, map[string]interface{}{"user": user, "token": token}, "Login successful", http.StatusOK)
}

func (h *AuthHandler) SetPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Token    string `json:"token"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendErrorResponse(w, "Invalid request", err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.authService.SetPassword(req.Token, req.Password); err != nil {
		h.sendErrorResponse(w, "Failed to set password", err.Error(), http.StatusBadRequest)
		return
	}

	h.sendSuccessResponse(w, nil, "Password set successfully. You can now login.", http.StatusOK)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err == nil {
		h.authService.Logout(cookie.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})

	h.sendSuccessResponse(w, nil, "Logged out successfully", http.StatusOK)
}

func (h *AuthHandler) sendErrorResponse(w http.ResponseWriter, error string, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Error: error, Message: message})
}

func (h *AuthHandler) sendSuccessResponse(w http.ResponseWriter, data interface{}, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(SuccessResponse{Success: true, Data: data, Message: message})
}
