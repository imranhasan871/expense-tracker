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
│       └── serve.go              # Server entry point & Route registration
├── internal/
│   ├── handlers/                 # HTTP Logic (API & UI)
│   │   ├── budget_handler.go
│   │   ├── category_handler.go
│   │   ├── expense_handler.go
│   │   └── template_handler.go
│   ├── models/                   # Domain Data Structures
│   │   ├── budget.go
│   │   ├── category.go
│   │   └── expense.go
│   └── repository/               # Database Access Layer
│       ├── budget_repository.go
│       ├── category_repository.go
│       └── expense_repository.go
├── migrations/                   # Auto-run SQL migrations
│   ├── 001_create_categories_table.sql
│   ├── 002_create_budgets_table.sql
│   ├── 003_create_expenses_table.sql
│   └── 004_seed_data.sql
├── web/
│   ├── static/                   # Assets (CSS & Interactivity)
│   │   ├── css/style.css
│   │   └── js/
│   │       ├── budgets.js
│   │       ├── categories.js
│   │       └── expenses.js
│   └── templates/                # HTML5 Components & Layouts
│       ├── index.html
│       ├── budgets.html
│       ├── categories.html
│       └── expenses.html
├── docker-compose.yml            # Multi-container orchestration
├── Dockerfile                    # Multi-stage optimized build
├── .dockerignore
├── .env                          # Configuration (DB URL, Port)
├── main.go                       # Application entry point
└── README.md
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
