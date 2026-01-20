package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type AdminHandler struct {
	db *sql.DB
}

func NewAdminHandler(db *sql.DB) *AdminHandler {
	return &AdminHandler{db: db}
}

func (h *AdminHandler) RunMigrations(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	authKey := r.Header.Get("X-Admin-Key")
	expectedKey := os.Getenv("ADMIN_SECRET_KEY")
	if expectedKey == "" {
		expectedKey = "your-secret-admin-key-change-this"
	}

	if authKey != expectedKey {
		http.Error(w, `{"error": "Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	migrationDir := "migrations"
	files, err := os.ReadDir(migrationDir)
	if err != nil {
		h.sendErrorResponse(w, "Failed to read migrations directory", err.Error())
		return
	}

	var sqlFiles []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".sql") {
			sqlFiles = append(sqlFiles, file.Name())
		}
	}
	sort.Strings(sqlFiles)

	results := []map[string]interface{}{}

	for _, filename := range sqlFiles {
		filePath := filepath.Join(migrationDir, filename)
		content, err := os.ReadFile(filePath)
		if err != nil {
			results = append(results, map[string]interface{}{
				"file":    filename,
				"status":  "error",
				"message": "Failed to read file: " + err.Error(),
			})
			continue
		}

		_, err = h.db.Exec(string(content))
		if err != nil {
			results = append(results, map[string]interface{}{
				"file":    filename,
				"status":  "error",
				"message": err.Error(),
			})
		} else {
			results = append(results, map[string]interface{}{
				"file":    filename,
				"status":  "success",
				"message": "Migration executed successfully",
			})
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Migration execution completed",
		"results": results,
	})
}

func (h *AdminHandler) sendErrorResponse(w http.ResponseWriter, error string, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(map[string]string{
		"error":   error,
		"message": message,
	})
}
