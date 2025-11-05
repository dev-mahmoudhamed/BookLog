CREATE TABLE books (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    author TEXT NOT NULL,
    published_year INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);