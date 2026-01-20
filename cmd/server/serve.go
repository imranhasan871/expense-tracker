package server

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"

	"expense-tracker/internal/handlers"
	"expense-tracker/internal/models"
	"expense-tracker/internal/repository"
	"expense-tracker/internal/service"
)

func Serve() {
	dbConnStr := os.Getenv("DATABASE_URL")
	if dbConnStr == "" {
		dbConnStr = "host=localhost port=5432 user=admin password=root dbname=expense_tracker sslmode=disable"
	}

	var db *sql.DB
	var err error

	// Retry connection - useful for slow database start
	for i := 0; i < 5; i++ {
		db, err = sql.Open("postgres", dbConnStr)
		if err == nil {
			if err = db.Ping(); err == nil {
				break
			}
		}
		log.Printf("Attempt %d: Failed to connect to database. Retrying in 2s...", i+1)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatal("Failed to connect to database after retries:", err)
	}
	defer db.Close()

	log.Println("✓ Connected to database!")

	// Run migrations
	if err := runMigrations(db); err != nil {
		log.Printf("Warning: Failed to run migrations: %v", err)
	}

	budgetRepo := repository.NewBudgetRepository(db)
	expenseRepo := repository.NewExpenseRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	userRepo := repository.NewUserRepository(db)

	emailService := service.NewEmailService()
	userService := service.NewUserService(userRepo, emailService)
	authService := service.NewAuthService(userRepo)

	categoryService := service.NewCategoryService(categoryRepo)
	if err := categoryService.InitializeDefaults(); err != nil {
		log.Printf("Warning: Failed to initialize default categories: %v", err)
	} else {
		log.Println("✓ Default categories initialized")
	}

	budgetService := service.NewBudgetService(budgetRepo, expenseRepo)
	expenseService := service.NewExpenseService(expenseRepo, budgetRepo)

	budgetHandler := handlers.NewBudgetHandler(budgetService)
	expenseHandler := handlers.NewExpenseHandler(expenseService)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	templateHandler := handlers.NewTemplateHandler("web/templates", categoryRepo, budgetRepo, expenseRepo)
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	authMiddleware := handlers.NewAuthMiddleware(authService)

	setupRoutes(categoryHandler, budgetHandler, expenseHandler, templateHandler, authHandler, userHandler, authMiddleware)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server running on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func runMigrations(db *sql.DB) error {
	migrationDir := "migrations"
	files, err := os.ReadDir(migrationDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}

		log.Printf("Running migration: %s", file.Name())
		content, err := os.ReadFile(migrationDir + "/" + file.Name())
		if err != nil {
			return err
		}

		_, err = db.Exec(string(content))
		if err != nil {
			return err
		}
	}
	log.Println("✓ All migrations completed")
	return nil
}

func setupRoutes(
	categoryHandler *handlers.CategoryHandler,
	budgetHandler *handlers.BudgetHandler,
	expenseHandler *handlers.ExpenseHandler,
	templateHandler *handlers.TemplateHandler,
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
	authMiddleware *handlers.AuthMiddleware,
) {
	http.HandleFunc("/", authMiddleware.Authenticate(templateHandler.RenderHome))
	http.HandleFunc("/categories", authMiddleware.RequireRole(models.RoleAdmin, models.RoleManagement)(templateHandler.RenderCategoriesPage))
	http.HandleFunc("/budgets", authMiddleware.RequireRole(models.RoleAdmin, models.RoleManagement)(templateHandler.RenderBudgetsPage))
	http.HandleFunc("/expenses", authMiddleware.RequireAuth(templateHandler.RenderExpensesPage))
	http.HandleFunc("/monitoring", authMiddleware.RequireRole(models.RoleAdmin, models.RoleManagement)(templateHandler.RenderMonitoringPage))
	http.HandleFunc("/login", authMiddleware.Authenticate(templateHandler.RenderLoginPage))
	http.HandleFunc("/set-password", authMiddleware.Authenticate(templateHandler.RenderSetPasswordPage))

	http.HandleFunc("/api/login", authHandler.Login)
	http.HandleFunc("/api/set-password", authHandler.SetPassword)
	http.HandleFunc("/api/logout", authHandler.Logout)

	// Admin only routes
	http.HandleFunc("/api/users", authMiddleware.RequireRole(models.RoleAdmin)(userHandler.ListUsers))
	http.HandleFunc("/api/users/create", authMiddleware.RequireRole(models.RoleAdmin)(userHandler.CreateUser))

	// Management routes
	http.HandleFunc("/api/categories", authMiddleware.RequireRole(models.RoleAdmin, models.RoleManagement)(categoryHandler.HandleCategories))
	http.HandleFunc("/api/categories/", authMiddleware.RequireRole(models.RoleAdmin, models.RoleManagement)(categoryHandler.HandleCategoryByID))

	http.HandleFunc("/api/budgets", authMiddleware.RequireRole(models.RoleAdmin, models.RoleManagement)(budgetHandler.HandleBudgets))
	http.HandleFunc("/api/budgets/status", authMiddleware.RequireRole(models.RoleAdmin, models.RoleManagement)(budgetHandler.GetBudgetStatus))
	http.HandleFunc("/api/monitoring", authMiddleware.RequireRole(models.RoleAdmin, models.RoleManagement)(budgetHandler.HandleMonitoring))

	http.HandleFunc("/api/budgets/", authMiddleware.RequireRole(models.RoleAdmin, models.RoleManagement)(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/lock") {
			budgetHandler.ToggleCircuitBreaker(w, r)
			return
		}
		http.NotFound(w, r)
	}))

	// Executive and above routes
	http.HandleFunc("/api/expenses", authMiddleware.Authenticate(expenseHandler.HandleExpenses))
	http.HandleFunc("/api/expenses/", authMiddleware.Authenticate(expenseHandler.HandleExpenseByID))

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
