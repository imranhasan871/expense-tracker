-- Create budgets table for tracking annual category limits
CREATE TABLE IF NOT EXISTS budgets (
    id SERIAL PRIMARY KEY,
    category_id INTEGER NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    amount DECIMAL(12, 2) NOT NULL DEFAULT 0.00,
    year INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(category_id, year)
);

-- Create index for faster lookups by category and year
CREATE INDEX IF NOT EXISTS idx_budgets_category_year ON budgets(category_id, year);

-- Add comment to table
COMMENT ON TABLE budgets IS 'Stores annual budget allocations for expense categories';
