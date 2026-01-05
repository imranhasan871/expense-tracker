-- Create expenses table for recording daily transactions
CREATE TABLE IF NOT EXISTS expenses (
    id SERIAL PRIMARY KEY,
    category_id INTEGER NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    amount DECIMAL(12, 2) NOT NULL DEFAULT 0.00,
    expense_date DATE NOT NULL DEFAULT CURRENT_DATE,
    remarks TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create index for faster filtering by date and category
CREATE INDEX IF NOT EXISTS idx_expenses_date ON expenses(expense_date);
CREATE INDEX IF NOT EXISTS idx_expenses_category ON expenses(category_id);

-- Add comment to table
COMMENT ON TABLE expenses IS 'Stores individual expense transactions linked to categories';
