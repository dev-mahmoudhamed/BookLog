package postgres

import (
	"bookLog/internal/models"
	"bookLog/internal/repository"
	"database/sql"
	"fmt"
)

type BookRepositoryPostgres struct {
	db *sql.DB
}

func NewBookRepositoryPostgres(db *sql.DB) repository.BookRepository {
	return &BookRepositoryPostgres{db: db}
}

func (r *BookRepositoryPostgres) GetAll(limit, offset int, search string) ([]models.Book, int, error) {
	query := `
        SELECT id, title, author, published_year, created_at
        FROM books
        WHERE 1=1
    `
	args := []interface{}{}
	idx := 1

	// --- 2. Add search filter ---
	if search != "" {
		query += fmt.Sprintf(" AND title ILIKE $%d", idx)
		args = append(args, "%"+search+"%")
		idx++
	}

	// --- 3. Add ordering, pagination ---
	query += fmt.Sprintf(" ORDER BY id ASC LIMIT $%d OFFSET $%d", idx, idx+1)
	args = append(args, limit, offset)

	// --- 4. Query rows ---
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var b models.Book

		if err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.PublishedYear, &b.CreatedAt); err != nil {
			return nil, 0, err
		}
		books = append(books, b)
	}

	// --- 5. Count total ---
	totalQuery := `SELECT COUNT(*) FROM books WHERE 1=1`
	totalArgs := []interface{}{}
	idx = 1

	if search != "" {
		totalQuery += fmt.Sprintf(" AND title ILIKE $%d", idx)
		totalArgs = append(totalArgs, "%"+search+"%")
	}

	var total int
	if err := r.db.QueryRow(totalQuery, totalArgs...).Scan(&total); err != nil {
		return nil, 0, err
	}

	return books, total, nil

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
