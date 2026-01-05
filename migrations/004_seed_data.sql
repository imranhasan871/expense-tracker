-- Seed initial categories (matching application defaults for consistency)
INSERT INTO categories (name, is_active) VALUES 
('Food', true),
('Transport', true),
('Rent', true),
('Utilities', true),
('Marketing', true),
('Salary', true),
('Office Rent', true),
('HR Development', true),
('Entertainment', true)
ON CONFLICT (name) DO NOTHING;

-- Seed some annual budgets for 2026
INSERT INTO budgets (category_id, amount, year)
SELECT id, 5000.00, 2026 FROM categories WHERE name = 'Food'
ON CONFLICT (category_id, year) DO NOTHING;

INSERT INTO budgets (category_id, amount, year)
SELECT id, 2000.00, 2026 FROM categories WHERE name = 'Transport'
ON CONFLICT (category_id, year) DO NOTHING;

INSERT INTO budgets (category_id, amount, year)
SELECT id, 12000.00, 2026 FROM categories WHERE name = 'Marketing'
ON CONFLICT (category_id, year) DO NOTHING;

INSERT INTO budgets (category_id, amount, year)
SELECT id, 1500.00, 2026 FROM categories WHERE name = 'Entertainment'
ON CONFLICT (category_id, year) DO NOTHING;

-- Seed some sample expenses
INSERT INTO expenses (category_id, amount, expense_date, remarks)
SELECT id, 84.50, '2026-01-05', 'Weekly Grocery Shopping at Big Bazaar' FROM categories WHERE name = 'Food';

INSERT INTO expenses (category_id, amount, expense_date, remarks)
SELECT id, 25.00, '2026-01-04', 'Movie Night with friends at Cineplex' FROM categories WHERE name = 'Entertainment';

INSERT INTO expenses (category_id, amount, expense_date, remarks)
SELECT id, 45.00, '2026-01-03', 'Uber ride to office' FROM categories WHERE name = 'Transport';

INSERT INTO expenses (category_id, amount, expense_date, remarks)
SELECT id, 120.00, '2026-01-02', 'Monthly Internet Bill' FROM categories WHERE name = 'Utilities';
