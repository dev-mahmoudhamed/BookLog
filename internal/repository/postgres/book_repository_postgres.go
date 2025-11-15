package postgres

import (
	"bookLog/internal/models"
	"bookLog/internal/repository"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
)

type BookRepositoryPostgres struct {
	db *sql.DB
}

func NewBookRepositoryPostgres(db *sql.DB) repository.BookRepository {
	return &BookRepositoryPostgres{db: db}
}

func (r *BookRepositoryPostgres) GetAll(limit, offset int, search string) ([]models.Book, int, error) {
	// Build the query with Squirrel
	query := sq.Select("id", "title", "author", "published_year", "created_at").
		From("books")

	// Add search filter if provided
	if search != "" {
		query = query.Where(sq.ILike{"title": "%" + search + "%"})
	}

	// Build and execute the query with pagination
	sql, args, err := query.
		OrderBy("id ASC").
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(sql, args...)
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

	// Count total records
	countQuery := sq.Select("COUNT(*)").From("books")

	if search != "" {
		countQuery = countQuery.Where(sq.ILike{"title": "%" + search + "%"})
	}

	countSql, countArgs, err := countQuery.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, 0, err
	}

	var total int
	if err := r.db.QueryRow(countSql, countArgs...).Scan(&total); err != nil {
		return nil, 0, err
	}

	return books, total, nil
}

func (r *BookRepositoryPostgres) GetByID(id int) (*models.Book, error) {
	sql, args, err := sq.Select("id", "title", "author", "published_year", "created_at").
		From("books").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, err
	}

	var b models.Book
	err = r.db.QueryRow(sql, args...).
		Scan(&b.ID, &b.Title, &b.Author, &b.PublishedYear, &b.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *BookRepositoryPostgres) Create(book *models.Book) error {
	sql, args, err := sq.Insert("books").
		Columns("title", "author", "published_year").
		Values(book.Title, book.Author, book.PublishedYear).
		Suffix("RETURNING id, created_at").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return err
	}

	err = r.db.QueryRow(sql, args...).Scan(&book.ID, &book.CreatedAt)
	return err
}

func (r *BookRepositoryPostgres) Update(book *models.Book) error {
	sql, args, err := sq.Update("books").
		Set("title", book.Title).
		Set("author", book.Author).
		Set("published_year", book.PublishedYear).
		Where(sq.Eq{"id": book.ID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return err
	}

	_, err = r.db.Exec(sql, args...)
	return err
}

func (r *BookRepositoryPostgres) Delete(id int) error {
	sql, args, err := sq.Delete("books").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return err
	}

	_, err = r.db.Exec(sql, args...)
	return err
}
