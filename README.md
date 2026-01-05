# Expense Tracker - Category Management

A Go-based expense tracker application with category management functionality. Built following Go project structure best practices with Go templates for the frontend.

## 📁 Project Structure

```
expense-tracker/
├── cmd/                          # Command-line applications
│   └── server/                   # Server application (optional)
├── internal/                     # Private application code
│   ├── handlers/                 # HTTP request handlers
│   │   ├── category_handler.go  # Category API handlers
│   │   └── template_handler.go  # Template rendering handlers
│   ├── models/                   # Data models
│   │   └── category.go          # Category model
│   └── repository/               # Data access layer
│       └── category_repository.go # Category database operations
├── migrations/                   # Database migrations
│   ├── 001_create_categories_table.sql
│   ├── 002_create_budgets_table.sql
│   └── 003_create_expenses_table.sql
├── web/                          # Web assets
│   ├── static/                   # Static files
│   └── templates/                # HTML templates
│       ├── index.html            # Home page
│       ├── categories.html       # Categories page
│       ├── budgets.html          # Budgets Planning page
│       └── expenses.html         # Expense Tracking page
├── Dockerfile                    # Docker build configuration
├── docker-compose.yml            # Docker orchestration
└── main.go                       # Application entry point
```

## 🚀 Features

### Core Modules
- ✅ **Category Management**: Organize expenses into meaningful groups.
- ✅ **Budget Planning**: Set annual limits for each category.
- ✅ **Expense Tracking**: Record daily transactions with remarks and filtering.

### Technical highlights
- ✅ **One-Command Setup**: Fully containerized with Docker.
- ✅ **Auto-Migrations**: Database schema initializes automatically on first run.
- ✅ **Sample Data**: Automatically seeds core categories, sample budgets, and expenses for a ready-to-use experience.

## 🛠️ Tech Stack

- **Backend**: Go 1.25+
- **Database**: PostgreSQL
- **Frontend**: Go Templates (html/template)
- **Styling**: Vanilla CSS with modern design
- **JavaScript**: Vanilla JS for interactivity

## 📦 Installation

The easiest way to run the project is using **Docker Compose**. This will start the application, the database with all necessary tables, and pgAdmin for database management.

### Method 1: Docker (Recommended)

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd expense-tracker
   ```

2. **Run with Docker Compose**
   ```bash
   docker-compose up -d
   ```

3. **Access the application**
   - **Web UI**: [http://localhost:8080](http://localhost:8080)
   - **pgAdmin**: [http://localhost:5050](http://localhost:5050) (Login: `admin@admin.com` / `root`)

### Method 2: Manual Setup (Local Development)

1. **Prerequisites**: Go 1.25+, PostgreSQL
2. **Start PostgreSQL** (standalone) and run migrations:
   ```bash
   psql -h localhost -U postgres -d expense_tracker -f migrations/001_create_categories_table.sql
   ```
3. **Download dependencies**
   ```bash
   go mod download
   ```
4. **Run the application**
   ```bash
   go run main.go
   ```


## 🌐 API Endpoints

### Web Routes (HTML)

| Method | Endpoint       | Description           |
|--------|----------------|-----------------------|
| GET    | `/`            | Home page             |
| GET    | `/categories`  | Categories management |

### API Routes (JSON)

| Method | Endpoint                | Description              |
|--------|-------------------------|--------------------------|
| GET    | `/api/categories`       | List all categories      |
| GET    | `/api/categories?active_only=true` | List active categories only |
| POST   | `/api/categories`       | Create new category      |
| GET    | `/api/categories/{id}`  | Get category by ID       |

## 📝 API Examples

### Create Category

**Request:**
```bash
curl -X POST http://localhost:8080/api/categories \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Travel",
    "is_active": true
  }'
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 10,
    "name": "Travel",
    "is_active": true,
    "created_at": "2026-01-05T15:41:26Z",
    "updated_at": "2026-01-05T15:41:26Z"
  },
  "message": "Category created successfully"
}
```

### Get All Categories

**Request:**
```bash
curl http://localhost:8080/api/categories
```

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "Food",
      "is_active": true,
      "created_at": "2026-01-05T10:00:00Z",
      "updated_at": "2026-01-05T10:00:00Z"
    },
    {
      "id": 2,
      "name": "Transport",
      "is_active": true,
      "created_at": "2026-01-05T10:00:00Z",
      "updated_at": "2026-01-05T10:00:00Z"
    }
  ]
}
```

### Get Category by ID

**Request:**
```bash
curl http://localhost:8080/api/categories/1
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "Food",
    "is_active": true,
    "created_at": "2026-01-05T10:00:00Z",
    "updated_at": "2026-01-05T10:00:00Z"
  }
}
```

## 🗄️ Database Schema

```sql
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

## 🎨 Design Features

- **Modern UI**: Gradient backgrounds, glassmorphism effects
- **Responsive**: Mobile-friendly design
- **Smooth Animations**: Hover effects and transitions
- **Premium Look**: Professional color scheme and typography
- **Accessible**: Semantic HTML and proper form labels

## 🔧 Configuration

Environment variables can be set in `.env` file:

```env
DATABASE_URL=host=localhost port=5432 user=postgres password=postgres dbname=expense sslmode=disable
PORT=8080
```

## 📚 Architecture

The project follows **Clean Architecture** principles:

- **Handlers**: Handle HTTP requests/responses
- **Repository**: Data access layer (database operations)
- **Models**: Domain entities
- **Separation of Concerns**: Clear boundaries between layers

## 🧪 Testing

Run tests:
```bash
go test ./...
```

## 📄 License

MIT License

## 👨‍💻 Author

Built with ❤️ using Go
