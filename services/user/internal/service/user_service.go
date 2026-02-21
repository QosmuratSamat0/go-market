package service

import (
	"context"
	"errors"

	userErr "github.com/go-market/pkg/errs"
	user "github.com/go-market/services/user/internal/model"
	userRepo "github.com/go-market/services/user/internal/repository"
)

type Service struct {
	repo userRepo.Repository
}

func New(repo userRepo.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetMe(ctx context.Context) (*user.User, error) {
	user, err := s.repo.GetMe(ctx)
	if err != nil {
		if errors.Is(err, userErr.ErrUserNotFound) {
			return nil, err
		}
		return nil, userErr.ErrFailedToGetUser
	}

	return user, err
}

func (s *Service) GetByID(ctx context.Context, id string) (*user.User, error) {
	if id == "" {
		return nil, userErr.ErrInvalidID
	}
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, userErr.ErrUserNotFound) {
			return nil, err
		}
		return nil, userErr.ErrFailedToGetUser
	}

	return user, err
}

func (s *Service) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	if email == "" {
		return nil, userErr.ErrInvalidEmail
	}
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, userErr.ErrUserNotFound) {
			return nil, err
		}
		return nil, userErr.ErrFailedToGetUser
	}

	return user, err
}

func (s *Service) Update(ctx context.Context, user user.User) error {
	if user.ID == "" {
		return userErr.ErrInvalidID
	}
	existingUser, err := s.repo.GetByID(ctx, user.ID)
	if err != nil {
		if errors.Is(err, userErr.ErrUserNotFound) {
			return err
		}
		return err
	}
	if existingUser == nil {
		return userErr.ErrUserNotFound
	}

	return s.repo.Update(ctx, user)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	if id == "" {
		return userErr.ErrInvalidID
	}

	existingUser, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, userErr.ErrUserNotFound) {
			return err
		}
	}
	if existingUser == nil {
		return userErr.ErrUserNotFound
	}

	return s.repo.Delete(ctx, id)
}
