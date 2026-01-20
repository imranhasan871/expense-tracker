package handlers

import (
	"context"
	"expense-tracker/internal/models"
	"expense-tracker/internal/service"
	"net/http"
	"strings"
)

type contextKey string

const (
	UserContextKey contextKey = "user"
)

type AuthMiddleware struct {
	authService service.AuthService
}

func NewAuthMiddleware(authService service.AuthService) *AuthMiddleware {
	return &AuthMiddleware{authService: authService}
}

func (m *AuthMiddleware) Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err == nil {
			user, authenticated := m.authService.IsAuthenticated(cookie.Value)
			if authenticated {
				ctx := context.WithValue(r.Context(), UserContextKey, user)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		}
		next.ServeHTTP(w, r)
	}
}

func (m *AuthMiddleware) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			m.handleUnauthorized(w, r)
			return
		}

		user, authenticated := m.authService.IsAuthenticated(cookie.Value)
		if !authenticated {
			m.handleUnauthorized(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func (m *AuthMiddleware) RequireRole(roles ...models.UserRole) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return m.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
			user := r.Context().Value(UserContextKey).(*models.User)

			allowed := false
			for _, role := range roles {
				if user.Role == role {
					allowed = true
					break
				}
			}

			if !allowed {
				if strings.HasPrefix(r.URL.Path, "/api/") {
					http.Error(w, `{"message": "Forbidden: insufficient permissions"}`, http.StatusForbidden)
				} else {
					http.Redirect(w, r, "/?error=forbidden", http.StatusSeeOther)
				}
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (m *AuthMiddleware) handleUnauthorized(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/api/") {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"message": "Unauthorized"}`))
		return
	}
	http.Redirect(w, r, "/login?return_to="+r.URL.Path, http.StatusSeeOther)
}

func GetAuthenticatedUser(r *http.Request) *models.User {
	user, _ := r.Context().Value(UserContextKey).(*models.User)
	return user
}
