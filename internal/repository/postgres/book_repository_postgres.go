package postgres

import (
	"bookLog/internal/models"
	"bookLog/internal/repository"
	"database/sql"
)

type BookRepositoryPostgres struct {
	db *sql.DB
}

func NewBookRepositoryPostgres(db *sql.DB) repository.BookRepository {
	return &BookRepositoryPostgres{db: db}
}

func (r *BookRepositoryPostgres) GetAll() ([]models.Book, error) {
	rows, err := r.db.Query("SELECT id, title, author, published_year, created_at FROM books ORDER BY id ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var b models.Book
		if err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.PublishedYear, &b.CreatedAt); err != nil {
			return nil, err
		}
		books = append(books, b)
	}
	return books, nil
}

func (r *BookRepositoryPostgres) GetByID(id int) (*models.Book, error) {
	var b models.Book
	err := r.db.QueryRow("SELECT id, title, author, published_year, created_at FROM books WHERE id=$1", id).
		Scan(&b.ID, &b.Title, &b.Author, &b.PublishedYear, &b.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *BookRepositoryPostgres) Create(book *models.Book) error {
	err := r.db.QueryRow("INSERT INTO books (title, author, published_year) VALUES ($1, $2, $3) RETURNING id, created_at",
		book.Title, book.Author, book.PublishedYear).Scan(&book.ID, &book.CreatedAt)
	return err
}

func (r *BookRepositoryPostgres) Update(book *models.Book) error {
	_, err := r.db.Exec("UPDATE books SET title = $1, author = $2, published_year = $3 WHERE id = $4",
		book.Title, book.Author, book.PublishedYear, book.ID)
	return err
}

func (r *BookRepositoryPostgres) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM books WHERE id=$1", id)
	return err
}
