package postgres

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	userErr "github.com/go-market/pkg/errs"
	user "github.com/go-market/services/user/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepo struct {
	db *pgxpool.Pool
}

func New(db_url string) (*PostgresRepo, error) {
	ctx := context.Background()
	db, err := pgxpool.New(ctx, db_url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect database %w", err)
	}

	return &PostgresRepo{db: db}, nil
}

func (r *PostgresRepo) Create(ctx context.Context, user user.User) (string, error) {
	const op = "repo.Create"
	slog.With("op", op)
	var id string
	query := `INSERT INTO users(name, email, avatar) VALUES ($1, $2, $3) RETURNING id`
	err := r.db.QueryRow(ctx, query, user.Username, user.Email, user.Avatar).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return "", fmt.Errorf("%s: %w", op, userErr.ErrUserExists)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (r *PostgresRepo) GetByID(ctx context.Context, id string) (*user.User, error) {
	const op = "repo.GetByID"
	slog.With("op", op)

	u := &user.User{}
	query := `SELECT id, name, email, avatar, created_at, updated_at FROM users WHERE id = $1`
	err := r.db.QueryRow(ctx, query, id).Scan(&u.ID, &u.Username, &u.Email, &u.Avatar, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, userErr.ErrUserNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return u, nil
}

func (r *PostgresRepo) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	const op = "repo.GetByEmail"
	slog.With("op", op)

	u := &user.User{}
	query := `SELECT id, name, email, avatar, created_at, updated_at FROM users WHERE email = $1`
	err := r.db.QueryRow(ctx, query, email).Scan(&u.ID, &u.Username, &u.Email, &u.Avatar, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return u, nil
}

func (r *PostgresRepo) Update(ctx context.Context, user user.User) error {
	const op = "repo.Update"
	slog.With("op", op)

	query := `UPDATE users SET name = $1, email = $2, avatar = $3 WHERE id = $4`
	result, err := r.db.Exec(ctx, query, user.Username, user.Email, user.Avatar, user.ID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("%s: user not found", op)
	}

	return nil
}

func (r *PostgresRepo) Delete(ctx context.Context, id string) error {
	const op = "repo.Delete"
	slog.With("op", op)

	query := `DELETE FROM users WHERE id = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("%s: user not found", op)
	}

	return nil
}

func (r *PostgresRepo) Close() {
	r.db.Close()
}
