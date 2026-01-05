# 💰 Premium Expense Tracker

A professional, full-stack financial management application built with **Go**, **PostgreSQL**, and **Docker**. This project features a modern, high-end UI with real-time data management for categories, budgets, and transactions.

## 🖼️ UI Previews

### Dashboard & Home
![Home Page](docs/screenshots/home.png)

### Expense Categories
![Categories Page](docs/screenshots/categories.png)

### Budget Planning
![Budgets Page](docs/screenshots/budgets.png)

### Transaction Tracking
![Expenses Page](docs/screenshots/expenses.png)

## 📁 Project Structure

```
expense-tracker/
├── cmd/
│   └── server/
│       └── serve.go              # Core server setup: dependency injection, router config, and middleware
├── internal/
│   ├── handlers/                 # Presentation Layer: Maps HTTP requests to repository logic
│   │   ├── budget_handler.go     # JSON API for budget CRUD and dashboard summaries
│   │   ├── category_handler.go   # API for category management (create, update, toggle status)
│   │   ├── expense_handler.go    # API for transaction management and advanced filtering
│   │   └── template_handler.go   # Server-Side Rendering: Prepares data for HTML templates
│   ├── models/                   # Domain Layer: Plain Go objects representing our business data
│   │   ├── budget.go             # Data structures for financial planning and annual targets
│   │   ├── category.go           # Data structures for expense classification
│   │   └── expense.go            # Data structures for individual financial records
│   └── repository/               # Data Access Layer: Isolated SQL logic using PostgreSQL
│       ├── budget_repository.go  # Optimized queries for annual budgets and upsert logic
│       ├── category_repository.go # Management of categories with duplicate-name protection
│       └── expense_repository.go # Dynamic query builder for date and category-based filtering
├── migrations/                   # Database Evolution: SQL scripts executed automatically on startup
│   ├── 001_create_categories_table.sql
│   ├── 002_create_budgets_table.sql
│   ├── 003_create_expenses_table.sql
│   └── 004_seed_data.sql        # Pre-fills the app with sample data for immediate demo
├── web/
│   ├── static/                   # Public Assets: Client-side logic and styling
│   │   ├── css/style.css        # Modern design system (Glassmorphism, Azure theme, Pill badges)
│   │   └── js/
│   │       ├── budgets.js       # AJAX logic for real-time budget updates and dashboard stats
│   │       ├── categories.js    # Logic for dynamic status toggling and local row filtering
│   │       └── expenses.js      # Transaction management and asynchronous list filtering
│   └── templates/                # View Layer: Modular HTML5 templates using Go's html/template
│       ├── index.html           # Landing page with feature overview
│       ├── budgets.html         # Interactive budget planning dashboard
│       ├── categories.html      # Category management interface with status toggles
│       └── expenses.html        # Transaction ledger with search and filters
├── docker-compose.yml            # Container Orchestration: Links Go app, PostgreSQL, and pgAdmin
├── Dockerfile                    # Multi-stage build for a lightweight, secure production image
├── .dockerignore                 # Excludes local files from Docker context to optimize builds
├── .env                          # Local configuration for database secrets and server ports
├── main.go                       # Minimal entry point that boots the cmd/server package
└── README.md                     # Project documentation and developer guide
```

## 🚀 Features

### 📊 Comprehensive Management
- **Category Control**: Create, update, and toggle active status for expense groups. Includes local filtering for active-only views.
- **Budget Intelligence**: Set annual targets per category with live dashboard summaries (Total Budget, Highest Allocation, Savings Target).
- **Transaction Ledger**: record daily expenses with remarks and dynamic filtering (date range, category).

### ⚡ Technical Excellence
- **Dockerized Architecture**: One-command deployment with Go, PostgreSQL, and pgAdmin.
- **Automated Schema**: Intelligent migrations that run on startup to prepare your database.
- **Transactional Integrity**: Robust repository layer with parameterized queries to prevent SQL injection.
- **Premium UX**: Modern Glassmorphism UI, semantic HTML5, and responsive Vanilla CSS.

## 📦 Installation

The application is designed to be up and running in seconds.

### Quick Start (Docker)

1. **Clone & Enter**
   ```bash
   git clone <repository-url>
   cd expense-tracker
   ```

2. **Launch Services**
   ```bash
   docker-compose up --build -d
   ```

3. **Access**
   - **Application**: [http://localhost:8080](http://localhost:8080)
   - **Database (pgAdmin)**: [http://localhost:5050](http://localhost:5050)
     - *User*: `admin@admin.com`
     - *Pass*: `root`

## 🌐 API Reference

### Categories (`/api/categories`)
- `GET /api/categories`: Fetch all categories
- `POST /api/categories`: Create category
- `PUT /api/categories/{id}`: Update category
- `PATCH /api/categories/{id}`: Toggle status (Active/Inactive)

### Budgets (`/api/budgets`)
- `GET /api/budgets?year=2026`: Fetch budgets and summary for a year
- `POST /api/budgets`: Set/Update budget for a category

### Expenses (`/api/expenses`)
- `GET /api/expenses`: List expenses with filters (`start_date`, `end_date`, `category_id`)
- `POST /api/expenses`: Record new transaction
- `DELETE /api/expenses/{id}`: Remove record

## �️ Configuration

Configure your environment in `.env`:
```env
DATABASE_URL=host=db port=5432 user=postgres password=postgres dbname=expense sslmode=disable
PORT=8080
```

---
Built with ❤️ and **Clean Architecture**.
