-- Create index for faster filtering by user_id in expenses table
CREATE INDEX IF NOT EXISTS idx_expenses_user_id ON expenses(user_id);
