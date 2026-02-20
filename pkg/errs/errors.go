package errs

import "errors"

var (
	// user
	ErrUserNotFound    = errors.New("user not found")
	ErrUserExists      = errors.New("user already exists")
	ErrInvalidID       = errors.New("invalid user id")
	ErrInvalidEmail    = errors.New("invalid user email")
	ErrFailedToGetUser = errors.New("failed to get user")
)
