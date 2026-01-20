package handlers

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"expense-tracker/internal/repository"
)

type TemplateHandler struct {
	templates   *template.Template
	catRepo     repository.CategoryRepository
	budgetRepo  repository.BudgetRepository
	expenseRepo repository.ExpenseRepository
}

func NewTemplateHandler(templatesDir string, catRepo repository.CategoryRepository, budgetRepo repository.BudgetRepository, expenseRepo repository.ExpenseRepository) *TemplateHandler {
	templates := template.Must(template.ParseGlob(filepath.Join(templatesDir, "*.html")))

	return &TemplateHandler{
		templates:   templates,
		catRepo:     catRepo,
		budgetRepo:  budgetRepo,
		expenseRepo: expenseRepo,
	}
}

func (h *TemplateHandler) RenderHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	data := struct {
		User  interface{}
		Title string
	}{
		User:  GetAuthenticatedUser(r),
		Title: "Home",
	}

	err := h.templates.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		log.Printf("Error rendering template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (h *TemplateHandler) RenderCategoriesPage(w http.ResponseWriter, r *http.Request) {
	categories, err := h.catRepo.GetAll(false)
	if err != nil {
		log.Printf("Error fetching categories: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Categories interface{}
		Title      string
		User       interface{}
	}{
		Categories: categories,
		Title:      "Expense Categories",
		User:       GetAuthenticatedUser(r),
	}

	err = h.templates.ExecuteTemplate(w, "categories.html", data)
	if err != nil {
		log.Printf("Error rendering template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (h *TemplateHandler) RenderBudgetsPage(w http.ResponseWriter, r *http.Request) {
	categories, err := h.catRepo.GetAll(true)
	if err != nil {
		log.Printf("Error fetching categories: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	summary, _ := h.budgetRepo.GetDashboardSummary(2026)

	data := struct {
		Categories interface{}
		Summary    interface{}
		Title      string
		User       interface{}
	}{
		Categories: categories,
		Summary:    summary,
		Title:      "Budget Planning",
		User:       GetAuthenticatedUser(r),
	}

	err = h.templates.ExecuteTemplate(w, "budgets.html", data)
	if err != nil {
		log.Printf("Error rendering template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (h *TemplateHandler) RenderExpensesPage(w http.ResponseWriter, r *http.Request) {
	categories, err := h.catRepo.GetAll(true)
	if err != nil {
		log.Printf("Error fetching categories: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Categories interface{}
		Title      string
		User       interface{}
	}{
		Categories: categories,
		Title:      "Expense Tracking",
		User:       GetAuthenticatedUser(r),
	}

	err = h.templates.ExecuteTemplate(w, "expenses.html", data)
	if err != nil {
		log.Printf("Error rendering template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (h *TemplateHandler) RenderMonitoringPage(w http.ResponseWriter, r *http.Request) {
	data := struct {
		User  interface{}
		Title string
	}{
		User:  GetAuthenticatedUser(r),
		Title: "Monitoring",
	}
	err := h.templates.ExecuteTemplate(w, "monitoring.html", data)
	if err != nil {
		log.Printf("Error rendering template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (h *TemplateHandler) RenderLoginPage(w http.ResponseWriter, r *http.Request) {
	data := struct {
		User  interface{}
		Title string
	}{
		User:  GetAuthenticatedUser(r),
		Title: "Login",
	}
	err := h.templates.ExecuteTemplate(w, "login.html", data)
	if err != nil {
		log.Printf("Error rendering template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (h *TemplateHandler) RenderSetPasswordPage(w http.ResponseWriter, r *http.Request) {
	data := struct {
		User  interface{}
		Title string
	}{
		User:  GetAuthenticatedUser(r),
		Title: "Set Password",
	}
	err := h.templates.ExecuteTemplate(w, "set-password.html", data)
	if err != nil {
		log.Printf("Error rendering template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (h *TemplateHandler) RenderUsersPage(w http.ResponseWriter, r *http.Request) {
	data := struct {
		User  interface{}
		Title string
	}{
		User:  GetAuthenticatedUser(r),
		Title: "User Management",
	}
	err := h.templates.ExecuteTemplate(w, "users.html", data)
	if err != nil {
		log.Printf("Error rendering template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
