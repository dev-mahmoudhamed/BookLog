package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"` // UUID primary key
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`      // Unique email
	Password  string    `json:"-"`          // Hashed password, never expose
	Role      string    `json:"role"`       // e.g. "admin", "user"
	CreatedAt time.Time `json:"created_at"` // Record creation timestamp
	UpdatedAt time.Time `json:"updated_at"` // Optional update timestamp
	DeletedAt time.Time `json:"deleted_at"` // Optional delete timestamp
}
