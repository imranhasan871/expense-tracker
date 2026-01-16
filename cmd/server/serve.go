package server

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strings"

	_ "github.com/lib/pq"

	"expense-tracker/internal/handlers"
	"expense-tracker/internal/repository"
)

func Serve() {
	dbConnStr := os.Getenv("DATABASE_URL")
	if dbConnStr == "" {
		dbConnStr = "host=localhost port=5432 user=admin password=root dbname=expense_tracker sslmode=disable"
	}

	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}
	log.Println("✓ Connected to database!")

	budgetRepo := repository.NewBudgetRepository(db)
	expenseRepo := repository.NewExpenseRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)

	if err := categoryRepo.InitializeDefaults(); err != nil {
		log.Printf("Warning: Failed to initialize default categories: %v", err)
	} else {
		log.Println("✓ Default categories initialized")
	}

	budgetHandler := handlers.NewBudgetHandler(budgetRepo, expenseRepo)
	expenseHandler := handlers.NewExpenseHandler(expenseRepo, budgetRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryRepo)
	templateHandler := handlers.NewTemplateHandler("web/templates", categoryRepo, budgetRepo, expenseRepo)

	setupRoutes(categoryHandler, budgetHandler, expenseHandler, templateHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server running on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func setupRoutes(
	categoryHandler *handlers.CategoryHandler,
	budgetHandler *handlers.BudgetHandler,
	expenseHandler *handlers.ExpenseHandler,
	templateHandler *handlers.TemplateHandler,
) {
	http.HandleFunc("/", templateHandler.RenderHome)
	http.HandleFunc("/categories", templateHandler.RenderCategoriesPage)
	http.HandleFunc("/budgets", templateHandler.RenderBudgetsPage)
	http.HandleFunc("/expenses", templateHandler.RenderExpensesPage)
	http.HandleFunc("/monitoring", templateHandler.RenderMonitoringPage)

	http.HandleFunc("/api/categories", categoryHandler.HandleCategories)
	http.HandleFunc("/api/categories/", categoryHandler.HandleCategoryByID)

	http.HandleFunc("/api/budgets", budgetHandler.HandleBudgets)
	http.HandleFunc("/api/budgets/status", budgetHandler.GetBudgetStatus)
	http.HandleFunc("/api/monitoring", budgetHandler.HandleMonitoring)

	http.HandleFunc("/api/budgets/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/lock") {
			budgetHandler.ToggleCircuitBreaker(w, r)
			return
		}
		http.NotFound(w, r)
	})

	http.HandleFunc("/api/expenses", expenseHandler.HandleExpenses)
	http.HandleFunc("/api/expenses/", expenseHandler.HandleExpenseByID)

	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Println("✓ Routes configured:")
	log.Println("  - GET  /                      (Home page)")
	log.Println("  - GET  /categories            (Categories page)")
	log.Println("  - GET  /budgets               (Budgets page)")
	log.Println("  - GET  /expenses              (Expenses page)")
	log.Println("  - GET  /monitoring            (Monitoring page)")
	log.Println("  - API  /api/categories        (Full CRUD)")
	log.Println("  - API  /api/budgets           (Set & Get)")
	log.Println("  - API  /api/expenses          (Filter & Record)")
	log.Println("  - GET  /static/*              (Static files)")
}
