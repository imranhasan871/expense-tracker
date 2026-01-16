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

	categoryRepo := repository.NewCategoryRepository(db)
	budgetRepo := repository.NewBudgetRepository(db)
	expenseRepo := repository.NewExpenseRepository(db)

	if err := categoryRepo.InitializeDefaults(); err != nil {
		log.Printf("Warning: Failed to initialize default categories: %v", err)
	} else {
		log.Println("✓ Default categories initialized")
	}

	categoryHandler := handlers.NewCategoryHandler(categoryRepo)
	budgetHandler := handlers.NewBudgetHandler(budgetRepo, expenseRepo)
	expenseHandler := handlers.NewExpenseHandler(expenseRepo, budgetRepo)
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

	// Dynamic route for lock - simple matching as Go 1.22+ mux is not used here (using stdlib 1.x style or need wrapper)
	// serve.go uses http.HandleFunc which is basic Mux in Go < 1.22.
	// For /api/budgets/{id}/lock, we can use a prefix or just match in a wrapper.
	// However, looking at handleCategoryByID, it seems we are just handling prefixes.
	// But `HandleMonitoring` handles `/api/monitoring`.
	// For `ToggleCircuitBreaker`, it processes `/api/budgets/.../lock`.
	// Let's rely on standard handling pattern shown in `HandleCategoryByID` (route /api/categories/).
	// We'll map `/api/budgets/` to a dispatcher if we want cleaner code, but `HandleBudgets` is mapped to `/api/budgets`.

	// To avoid conflicts with `/api/budgets` (exact match usually? no, HandleFunc is prefix based if ends with /)
	// /api/budgets is NOT ending with slash in previous code.

	// We need a way to route /api/budgets/{id}/lock.
	// Let's add a specific handler for this prefix.
	// Since HandleBudgets is on `/api/budgets`, it might catch others if we changed it to `/api/budgets/`.
	// Use a closure or specific path if possible.
	// Actually, the simplest way given the existing code structure is:

	http.HandleFunc("/api/budgets/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/lock") && r.Method == "POST" {
			budgetHandler.ToggleCircuitBreaker(w, r)
			return
		}
		// Fallback doesn't work well with http.HandleFunc logic if we want to support both.
		// The existing `/api/budgets` handler handles GET/POST for root.
		// If we register `/api/budgets/` (with slash), it handles all subpaths.
		// So we should move `HandleBudgets` logic here or separate them.

		// Strategy: Leave `/api/budgets` for the main handler.
		// Register `/api/budgets/` for the children.
		// BUT `http.ServeMux` prioritizes longest match.
		// So this works:
		// /api/budgets -> HandleBudgets
		// /api/budgets/ -> New Dispatcher
	})

	// Wait, to keep it simple and consistent with `HandleCategoryByID`:
	// We will register specifics first?? No.
	// Let's try to match how others are done.
	// `HandleCategoryByID` is at `/api/categories/`. `HandleCategories` is at `/api/categories`.

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
	log.Println("  - API  /api/categories        (Full CRUD)")
	log.Println("  - API  /api/budgets           (Set & Get)")
	log.Println("  - API  /api/expenses          (Filter & Record)")
	log.Println("  - GET  /static/*              (Static files)")
}
