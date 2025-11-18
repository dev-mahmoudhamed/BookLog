package repository

import (
	"authService/internal/models"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

type UserRepositoryPostgres struct {
	db *sql.DB
}

func NewUserRepositoryPostgres(db *sql.DB) UserRepository {
	return &UserRepositoryPostgres{db: db}
}

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
