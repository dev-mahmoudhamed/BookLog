package models

import "time"

type Book struct {
	ID            int       `json:"id"`
	Title         string    `json:"title"`
	Author        string    `json:"author"`
	PublishedYear int       `json:"published_year"`
	CreatedAt     time.Time `json:"created_at"`
}
