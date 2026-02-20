package model

import "time"

type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"name" validate:"required"`
	Email     string    `json:"email" validate:"required,email"`
	Avatar    string    `json:"avatar" validate:"omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
