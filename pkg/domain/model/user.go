package model

import "time"

type User struct {
	ID    string
	Name  string
	Email string
	Role  string
}

type UserRegisteredEvent struct {
	UserID       string    `json:"user_id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Role         string    `json:"role"`
	RegisteredAt time.Time `json:"registered_at"`
}
