package postgres

import (
	"bookLog/internal/models"
	"bookLog/internal/repository"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

type UserRepositoryPostgres struct {
	db *sql.DB
}

func NewUserRepositoryPostgres(db *sql.DB) repository.UserRepository {
	return &UserRepositoryPostgres{db: db}
}

// Register implements repository.UserRepository.
func (u *UserRepositoryPostgres) Create(user *models.User) error {
	// ensure ID and timestamps
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	now := time.Now().UTC()
	if user.CreatedAt.IsZero() {
		user.CreatedAt = now
	}
	user.UpdatedAt = now

	query := sq.Insert("users").
		Columns("id", "full_name", "email", "password", "role", "created_at", "updated_at").
		Values(user.ID, user.FullName, user.Email, user.Password, user.Role, user.CreatedAt, user.UpdatedAt).
		PlaceholderFormat(sq.Dollar)

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = u.db.Exec(sqlStr, args...)
	return err
}

// Delete implements repository.UserRepository.
func (u *UserRepositoryPostgres) Delete(id uuid.UUID) error {
	query := sq.Delete("users").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = u.db.Exec(sqlStr, args...)
	return err
}

// GetByID implements repository.UserRepository.
func (u *UserRepositoryPostgres) GetByID(id uuid.UUID) (*models.User, error) {
	query := sq.Select("id", "full_name", "email", "password", "role", "created_at", "updated_at").
		From("users").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		Limit(1)

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var usr models.User
	row := u.db.QueryRow(sqlStr, args...)
	if err := row.Scan(&usr.ID, &usr.FullName, &usr.Email, &usr.Password, &usr.Role, &usr.CreatedAt, &usr.UpdatedAt); err != nil {
		return nil, err
	}
	return &usr, nil
}

func (u *UserRepositoryPostgres) GetByEmail(email string) (*models.User, error) {
	query := sq.Select("id", "full_name", "email", "password", "role", "created_at", "updated_at").
		From("users").
		Where(sq.Eq{"email": email}).
		PlaceholderFormat(sq.Dollar).
		Limit(1)
	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var usr models.User
	row := u.db.QueryRow(sqlStr, args...)
	if err := row.Scan(&usr.ID, &usr.FullName, &usr.Email, &usr.Password, &usr.Role, &usr.CreatedAt, &usr.UpdatedAt); err != nil {
		return nil, err
	}
	return &usr, nil
}

// Update implements repository.UserRepository.
func (u *UserRepositoryPostgres) Update(user *models.User) error {
	user.UpdatedAt = time.Now().UTC()

	query := sq.Update("users").
		Set("fullname", user.FullName).
		Set("email", user.Email).
		Set("password", user.Password).
		Set("role", user.Role).
		Set("updated_at", user.UpdatedAt).
		Where(sq.Eq{"id": user.ID}).
		PlaceholderFormat(sq.Dollar)

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = u.db.Exec(sqlStr, args...)
	return err
}

func (r *UserRepositoryPostgres) GetAll(search string) ([]models.User, int, error) {
	query := sq.Select("id", "full_name", "email", "role", "created_at", "updated_at").From("users")

	if search != "" {
		query = query.Where(sq.ILike{"fullname": "%" + search + "%"})
	}

	sql, args, err := query.
		OrderBy("id ASC").
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

	var users []models.User
	for rows.Next() {
		var b models.User
		if err := rows.Scan(&b.ID, &b.FullName, &b.Email, &b.Role, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, 0, err
		}
		users = append(users, b)
	}

	countQuery := sq.Select("COUNT(*)").From("users")

	if search != "" {
		countQuery = countQuery.Where(sq.ILike{"fullname": "%" + search + "%"})
	}

	countSql, countArgs, err := countQuery.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, 0, err
	}

	var total int
	if err := r.db.QueryRow(countSql, countArgs...).Scan(&total); err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
