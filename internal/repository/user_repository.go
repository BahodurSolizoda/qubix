package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"qubix/internal/model"
	"qubix/internal/service"

	"github.com/lib/pq"
)

var _ service.UserRepository = (*UserRepository)(nil)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users (login, email, password_hash, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		user.Login,
		user.Email,
		user.PasswordHash,
		user.Role,
	).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return model.ErrDuplicate
		}

		return fmt.Errorf("create user: %w", err)
	}

	return nil
}

func (r *UserRepository) GetByLoginOrEmail(ctx context.Context, loginOrEmail string) (*model.User, error) {
	query := `
		SELECT id, login, email, password_hash, role, created_at, updated_at
		FROM users
		WHERE login = $1 OR email = $1
	`

	var user model.User

	err := r.db.QueryRowContext(ctx, query, loginOrEmail).Scan(
		&user.ID,
		&user.Login,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("get user by login or email: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int) (*model.User, error) {
	query := `
		SELECT id, login, email, password_hash, role, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user model.User

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Login,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
	query := `
        UPDATE users
        SET login = $1, email = $2, updated_at = CURRENT_TIMESTAMP
        WHERE id = $3
    `

	_, err := r.db.ExecContext(ctx, query, user.Login, user.Email, user.ID)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}

	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	return nil
}
