package models

import "time"

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name" binding:"required"`
	Number    string    `json:"number" binding:"required"`
	Email     string    `json:"email" binding:"required,email"`
	Password  string    `json:"-"` // Ne pas exposer dans JSON
	APIKey    string    `json:"api_key"`
	CreatedAt time.Time `json:"created_at"`
}
