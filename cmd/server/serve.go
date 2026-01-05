package server

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"expense-tracker/internal/handlers"
	"expense-tracker/internal/repository"
)

func Serve() {
	// Get database connection string from environment or use default
	dbConnStr := os.Getenv("DATABASE_URL")
	if dbConnStr == "" {
		dbConnStr = "host=localhost port=5432 user=admin password=root dbname=expense_tracker sslmode=disable"
	}

	// Connect to PostgreSQL
	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Test connection
	if err = db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}
	log.Println("✓ Connected to database!")

	// Initialize repository
	categoryRepo := repository.NewCategoryRepository(db)

	// Initialize default categories
	if err := categoryRepo.InitializeDefaults(); err != nil {
		log.Printf("Warning: Failed to initialize default categories: %v", err)
	} else {
		log.Println("✓ Default categories initialized")
	}

	// Initialize handlers
	categoryHandler := handlers.NewCategoryHandler(categoryRepo)
	templateHandler := handlers.NewTemplateHandler("web/templates", categoryRepo)

	// Setup routes
	setupRoutes(categoryHandler, templateHandler)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server running on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func setupRoutes(categoryHandler *handlers.CategoryHandler, templateHandler *handlers.TemplateHandler) {
	// Web routes (HTML templates)
	http.HandleFunc("/", templateHandler.RenderHome)
	http.HandleFunc("/categories", templateHandler.RenderCategoriesPage)
	http.HandleFunc("/budgets", templateHandler.RenderBudgetsPage)

	// API routes (JSON)
	http.HandleFunc("/api/categories", categoryHandler.HandleCategories)
	http.HandleFunc("/api/categories/", categoryHandler.HandleCategoryByID)

	// Static files
	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Println("✓ Routes configured:")
	log.Println("  - GET  /                      (Home page)")
	log.Println("  - GET  /categories            (Categories page)")
	log.Println("  - GET  /budgets               (Budgets page)")
	log.Println("  - GET  /api/categories        (List all categories)")
	log.Println("  - POST /api/categories        (Create category)")
	log.Println("  - GET  /api/categories/{id}   (Get category by ID)")
	log.Println("  - GET  /static/*              (Static files)")
}
