package repository

import (
	"context"

	user "github.com/go-market/services/user/internal/model"
)

type Repository interface {
	GetMe(ctx context.Context) (*user.User, error)
	GetByID(ctx context.Context, id string) (*user.User, error)
	GetByEmail(ctx context.Context, email string) (*user.User, error)
	Update(ctx context.Context, user user.User) error
	Delete(ctx context.Context, id string) error
}
