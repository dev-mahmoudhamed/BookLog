package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Book struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Title       string `gorm:"type:varchar(255);not null" json:"title"`
	Author      string `gorm:"type:varchar(255);not null" json:"author"`
	Description string `gorm:"type:text" json:"description"`
	Year        int    `gorm:"type:int" json:"year"`
	UserID      uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
