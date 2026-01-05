package handlers

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"expense-tracker/internal/repository"
)

// TemplateHandler handles rendering of HTML templates
type TemplateHandler struct {
	templates *template.Template
	repo      *repository.CategoryRepository
}

// NewTemplateHandler creates a new template handler
func NewTemplateHandler(templatesDir string, repo *repository.CategoryRepository) *TemplateHandler {
	// Parse all templates
	templates := template.Must(template.ParseGlob(filepath.Join(templatesDir, "*.html")))

	return &TemplateHandler{
		templates: templates,
		repo:      repo,
	}
}

// RenderHome renders the home page
func (h *TemplateHandler) RenderHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	err := h.templates.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		log.Printf("Error rendering template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// RenderCategoriesPage renders the categories management page
func (h *TemplateHandler) RenderCategoriesPage(w http.ResponseWriter, r *http.Request) {
	// Get all categories from database
	categories, err := h.repo.GetAll(false)
	if err != nil {
		log.Printf("Error fetching categories: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Categories interface{}
		Title      string
	}{
		Categories: categories,
		Title:      "Expense Categories",
	}

	err = h.templates.ExecuteTemplate(w, "categories.html", data)
	if err != nil {
		log.Printf("Error rendering template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// RenderBudgetsPage renders the budgets management page
func (h *TemplateHandler) RenderBudgetsPage(w http.ResponseWriter, r *http.Request) {
	// We need categories for the budget creation dropdown
	categories, err := h.repo.GetAll(true) // Fetch only active categories
	if err != nil {
		log.Printf("Error fetching categories: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Categories interface{}
		Title      string
	}{
		Categories: categories,
		Title:      "Budget Planning",
	}

	err = h.templates.ExecuteTemplate(w, "budgets.html", data)
	if err != nil {
		log.Printf("Error rendering template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// RenderExpensesPage renders the expenses management page
func (h *TemplateHandler) RenderExpensesPage(w http.ResponseWriter, r *http.Request) {
	// Need categories for the filter and creation form
	categories, err := h.repo.GetAll(true)
	if err != nil {
		log.Printf("Error fetching categories: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Categories interface{}
		Title      string
	}{
		Categories: categories,
		Title:      "Expense Tracking",
	}

	err = h.templates.ExecuteTemplate(w, "expenses.html", data)
	if err != nil {
		log.Printf("Error rendering template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
