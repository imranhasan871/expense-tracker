-- Create categories table for expense category management
CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create index on name for faster lookups
CREATE INDEX IF NOT EXISTS idx_categories_name ON categories(name);

-- Create index on is_active for filtering
CREATE INDEX IF NOT EXISTS idx_categories_is_active ON categories(is_active);

-- Add comment to table
COMMENT ON TABLE categories IS 'Stores expense categories for organizing expenses';
COMMENT ON COLUMN categories.name IS 'Unique category name (case-insensitive uniqueness enforced at application level)';
COMMENT ON COLUMN categories.is_active IS 'Indicates if category is active. Inactive categories preserve historical data';
