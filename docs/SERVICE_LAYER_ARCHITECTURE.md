# Service Layer Architecture

## Overview
The service layer has been successfully introduced to the expense tracker application, following clean architecture principles and separation of concerns.

## Architecture Layers

```
┌─────────────────────────────────────┐
│         HTTP Handlers               │  ← Presentation Layer
│  (category, budget, expense)        │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│         Service Layer               │  ← Business Logic Layer
│  • CategoryService                  │
│  • BudgetService                    │
│  • ExpenseService                   │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│       Repository Layer              │  ← Data Access Layer
│  • CategoryRepository               │
│  • BudgetRepository                 │
│  • ExpenseRepository                │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│         Database (PostgreSQL)       │
└─────────────────────────────────────┘
```

## Service Layer Components

### 1. CategoryService (`internal/service/category_service.go`)
**Responsibilities:**
- Business logic for category management
- Input validation (trim whitespace, check empty names)
- Default category initialization

**Key Methods:**
- `GetAll(activeOnly bool)` - Retrieve all categories
- `GetByID(id int)` - Get category by ID
- `Create(name string, isActive bool)` - Create new category with validation
- `Update(id int, name string, isActive bool)` - Update existing category
- `ToggleStatus(id int)` - Toggle category active status
- `InitializeDefaults()` - Initialize default categories

**Business Rules:**
- Category names must not be empty after trimming
- Duplicate checking is delegated to repository

### 2. BudgetService (`internal/service/budget_service.go`)
**Responsibilities:**
- Budget management business logic
- Budget status calculation
- Circuit breaker (lock) management
- Input validation

**Key Methods:**
- `GetAll(year int)` - Get all budgets for a year
- `GetDashboardSummary(year int)` - Get budget summary statistics
- `GetStatus(categoryID, year int)` - Calculate budget status with spent amount
- `CreateOrUpdate(categoryID int, amount float64, year int)` - Create or update budget
- `ToggleLock(budgetID int, isLocked bool)` - Toggle circuit breaker
- `GetMonitoringData(year int)` - Get monitoring statistics
- `IsLocked(categoryID, year int)` - Check if budget is locked

**Business Rules:**
- Year must be greater than 0
- Category ID must be greater than 0
- Amount must be greater than 0
- Budget status calculation includes:
  - Allocated amount
  - Spent amount (from expenses)
  - Remaining amount
  - Percentage used
  - Lock status

### 3. ExpenseService (`internal/service/expense_service.go`)
**Responsibilities:**
- Expense management business logic
- Circuit breaker enforcement
- Input validation

**Key Methods:**
- `Create(req models.ExpenseRequest)` - Create expense with circuit breaker check
- `GetAll(filter models.ExpenseFilter)` - Get filtered expenses
- `Delete(id int)` - Delete expense

**Business Rules:**
- Category ID must be greater than 0
- Amount must be greater than 0
- Expense date is required
- **Circuit Breaker**: Prevents expense creation if budget is locked
- Filter validation is performed before querying

## Key Benefits

### 1. **Separation of Concerns**
- Handlers focus on HTTP request/response handling
- Services contain business logic
- Repositories handle data access

### 2. **Testability**
- Services can be tested independently
- Mock repositories can be injected for unit testing
- Business logic is isolated from HTTP and database concerns

### 3. **Reusability**
- Business logic can be reused across different handlers
- Services can be used by CLI, gRPC, or other interfaces

### 4. **Maintainability**
- Business rules are centralized in services
- Changes to business logic don't affect handlers or repositories
- Clear boundaries between layers

### 5. **Validation**
- Input validation is centralized in services
- Consistent validation across all entry points
- Reduces duplication in handlers

## Dependency Injection

Services are created in `cmd/server/serve.go`:

```go
// Create repositories
budgetRepo := repository.NewBudgetRepository(db)
expenseRepo := repository.NewExpenseRepository(db)
categoryRepo := repository.NewCategoryRepository(db)

// Create services (injecting repositories)
categoryService := service.NewCategoryService(categoryRepo)
budgetService := service.NewBudgetService(budgetRepo, expenseRepo)
expenseService := service.NewExpenseService(expenseRepo, budgetRepo)

// Create handlers (injecting services)
categoryHandler := handlers.NewCategoryHandler(categoryService)
budgetHandler := handlers.NewBudgetHandler(budgetService)
expenseHandler := handlers.NewExpenseHandler(expenseService)
```

## Interface Segregation

Each service defines its own repository interface requirements:

- **CategoryService** requires: `CategoryRepository` interface
- **BudgetService** requires: `BudgetRepository` and `ExpenseRepository` interfaces
- **ExpenseService** requires: `ExpenseRepositoryInterface` and `BudgetRepositoryInterface`

This follows the **Dependency Inversion Principle** - services depend on abstractions (interfaces), not concrete implementations.

## Changes Made

### Files Created:
1. `internal/service/category_service.go` - Category business logic
2. `internal/service/budget_service.go` - Budget business logic
3. `internal/service/expense_service.go` - Expense business logic

### Files Modified:
1. `internal/handlers/category_handler.go` - Now uses CategoryService
2. `internal/handlers/budget_handler.go` - Now uses BudgetService
3. `internal/handlers/expense_handler.go` - Now uses ExpenseService
4. `internal/models/budget.go` - Added BudgetStatus struct
5. `cmd/server/serve.go` - Wired up services with dependency injection

### Files Deleted:
1. `internal/service/interfaces.go` - Removed duplicate interface definitions

## Future Enhancements

1. **Add Service Interfaces**: Define interfaces for services to enable mocking in handler tests
2. **Add Logging**: Implement structured logging in services
3. **Add Metrics**: Track service-level metrics (call counts, durations)
4. **Add Caching**: Implement caching layer in services
5. **Add Transaction Support**: Wrap multi-step operations in database transactions
6. **Add Event Publishing**: Emit events for important business actions

## Testing Strategy

### Unit Tests (Services):
```go
// Mock repository
type MockCategoryRepository struct {
    GetAllFunc func(activeOnly bool) ([]models.Category, error)
}

// Test service
func TestCategoryService_Create(t *testing.T) {
    mockRepo := &MockCategoryRepository{...}
    service := service.NewCategoryService(mockRepo)
    // Test business logic
}
```

### Integration Tests (Handlers):
- Test handlers with real services and mock repositories
- Verify HTTP request/response handling
- Validate error responses

## Conclusion

The service layer successfully decouples business logic from HTTP handling and data access, making the application more maintainable, testable, and scalable. The architecture now follows clean architecture principles with clear separation of concerns.
