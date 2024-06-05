CREATE TABLE todos (
  id SERIAL PRIMARY KEY, 
  deadline TIMESTAMP, 
  task TEXT NOT NULL, 
  is_completed BOOLEAN NOT NULL DEFAULT FALSE, 
  todo_parent INTEGER NULL, 
  version INTEGER NOT NULL,
  CONSTRAINT todos_todos_parent FOREIGN KEY (todo_parent) REFERENCES todos (id) ON DELETE CASCADE
);
