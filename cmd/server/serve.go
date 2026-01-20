package server

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
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
	adminHandler := handlers.NewAdminHandler(db)
	authMiddleware := handlers.NewAuthMiddleware(authService)

	setupRoutes(categoryHandler, budgetHandler, expenseHandler, templateHandler, authHandler, userHandler, adminHandler, authMiddleware)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server running on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func runMigrations(db *sql.DB) error {
	var tableExists bool
	err := db.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'users')").Scan(&tableExists)
	if err != nil {
		log.Printf("Warning: Could not check for users table: %v", err)
	}

	if !tableExists {
		log.Println("⚠️  Users table not found - running ALL migrations...")
	}

	migrationDir := "migrations"
	files, err := os.ReadDir(migrationDir)
	if err != nil {
		return err
	}

	var sqlFiles []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".sql") {
			sqlFiles = append(sqlFiles, file.Name())
		}
	}
	sort.Strings(sqlFiles)

	successCount := 0
	for _, filename := range sqlFiles {
		filePath := filepath.Join(migrationDir, filename)
		content, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("Warning: Failed to read %s: %v", filename, err)
			continue
		}

		log.Printf("Running migration: %s", filename)
		_, err = db.Exec(string(content))
		if err != nil {
			errMsg := err.Error()
			if strings.Contains(errMsg, "already exists") ||
				strings.Contains(errMsg, "duplicate") ||
				strings.Contains(errMsg, "ON CONFLICT") {
				log.Printf("  ✓ %s (already applied)", filename)
				successCount++
			} else {
				log.Printf("  ✗ Failed: %v", err)
			}
		} else {
			log.Printf("  ✓ %s", filename)
			successCount++
		}
	}

	if successCount > 0 {
		log.Printf("✓ Successfully applied %d migrations", successCount)
	}

	return nil
}

func setupRoutes(
	categoryHandler *handlers.CategoryHandler,
	budgetHandler *handlers.BudgetHandler,
	expenseHandler *handlers.ExpenseHandler,
	templateHandler *handlers.TemplateHandler,
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
	adminHandler *handlers.AdminHandler,
	authMiddleware *handlers.AuthMiddleware,
) {
	http.HandleFunc("/", authMiddleware.Authenticate(templateHandler.RenderHome))
	http.HandleFunc("/categories", authMiddleware.RequireRole(models.RoleAdmin, models.RoleManagement)(templateHandler.RenderCategoriesPage))
	http.HandleFunc("/budgets", authMiddleware.RequireRole(models.RoleAdmin, models.RoleManagement)(templateHandler.RenderBudgetsPage))
	http.HandleFunc("/expenses", authMiddleware.RequireAuth(templateHandler.RenderExpensesPage))
	http.HandleFunc("/monitoring", authMiddleware.RequireRole(models.RoleAdmin, models.RoleManagement)(templateHandler.RenderMonitoringPage))
	http.HandleFunc("/users", authMiddleware.RequireRole(models.RoleAdmin)(templateHandler.RenderUsersPage))
	http.HandleFunc("/login", authMiddleware.Authenticate(templateHandler.RenderLoginPage))
	http.HandleFunc("/set-password", authMiddleware.Authenticate(templateHandler.RenderSetPasswordPage))

	http.HandleFunc("/api/login", authHandler.Login)
	http.HandleFunc("/api/set-password", authHandler.SetPassword)
	http.HandleFunc("/api/logout", authHandler.Logout)

	http.HandleFunc("/api/users", authMiddleware.RequireRole(models.RoleAdmin)(userHandler.ListUsers))
	http.HandleFunc("/api/users/create", authMiddleware.RequireRole(models.RoleAdmin)(userHandler.CreateUser))
	http.HandleFunc("/api/users/update-role", authMiddleware.RequireRole(models.RoleAdmin)(userHandler.UpdateUserRole))
	http.HandleFunc("/admin/run-migrations", adminHandler.RunMigrations)

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

	http.HandleFunc("/api/expenses", authMiddleware.Authenticate(expenseHandler.HandleExpenses))
	http.HandleFunc("/api/expenses/", authMiddleware.Authenticate(expenseHandler.HandleExpenseByID))

	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=86400")
		fs.ServeHTTP(w, r)
	})))

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
