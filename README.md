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
│   └── 001_create_categories_table.sql
├── web/                          # Web assets
│   ├── static/                   # Static files
│   │   ├── css/
│   │   │   └── style.css        # Application styles
│   │   └── js/
│   │       └── categories.js    # Frontend JavaScript
│   └── templates/                # HTML templates
│       ├── index.html           # Home page
│       └── categories.html      # Categories management page
├── .env                          # Environment variables
├── docker-compose.yml            # Docker compose configuration
├── go.mod                        # Go module definition
├── go.sum                        # Go dependencies checksums
└── main.go                       # Application entry point
```

## 🚀 Features

### Expense Category Management

- ✅ **Default Categories**: Automatically initializes with predefined categories:
  - Food, Transport, Rent, Utilities
  - Marketing, Salary, Office Rent
  - HR Development, Entertainment

- ✅ **CRUD Operations**:
  - Create new categories
  - View all categories
  - Get category by ID
  - Filter active/inactive categories

- ✅ **Business Logic**:
  - Unique category names (case-insensitive)
  - Active/Inactive status management
  - Preserves historical data (inactive categories)

## 🛠️ Tech Stack

- **Backend**: Go 1.25+
- **Database**: PostgreSQL
- **Frontend**: Go Templates (html/template)
- **Styling**: Vanilla CSS with modern design
- **JavaScript**: Vanilla JS for interactivity

## 📦 Installation

### Prerequisites

- Go 1.25 or higher
- PostgreSQL
- Docker & Docker Compose (optional)

### Setup

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd expense-tracker
   ```

2. **Start PostgreSQL** (using Docker)
   ```bash
   docker-compose up -d
   ```

3. **Run database migrations**
   ```bash
   psql -h localhost -U postgres -d expense -f migrations/001_create_categories_table.sql
   ```

4. **Install dependencies**
   ```bash
   go mod download
   ```

5. **Run the application**
   ```bash
   go run main.go
   ```

6. **Access the application**
   - Web UI: http://localhost:8080
   - Categories Page: http://localhost:8080/categories

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
