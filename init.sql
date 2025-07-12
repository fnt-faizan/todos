-- Create todos table
CREATE TABLE IF NOT EXISTS todos (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    status BOOLEAN DEFAULT FALSE,
    deleted BOOLEAN DEFAULT FALSE
);

-- Create indexes for faster lookups
CREATE INDEX idx_todos_status ON todos(status);
CREATE INDEX idx_todos_deleted ON todos(deleted);
