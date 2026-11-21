package repository

import (
	"context"
	"database/sql"
	"errors"

	"qubix/internal/model"
)

type AttemptRepository struct {
	db *sql.DB
}

func NewAttemptRepository(db *sql.DB) *AttemptRepository {
	return &AttemptRepository{
		db: db,
	}
}

func (r *AttemptRepository) Create(ctx context.Context, attempt *model.Attempt) error {
	query := `
		INSERT INTO attempts (user_id, quiz_id, status)
		VALUES ($1, $2, $3)
		RETURNING id, started_at
	`

	return r.db.QueryRowContext(
		ctx,
		query,
		attempt.UserID,
		attempt.QuizID,
		attempt.Status,
	).Scan(&attempt.ID, &attempt.StartedAt)
}

func (r *AttemptRepository) GetByID(ctx context.Context, id int) (*model.Attempt, error) {
	attempt := &model.Attempt{}

	query := `
		SELECT id, user_id, quiz_id, started_at, completed_at, status
		FROM attempts
		WHERE id = $1
	`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&attempt.ID,
		&attempt.UserID,
		&attempt.QuizID,
		&attempt.StartedAt,
		&attempt.CompletedAt,
		&attempt.Status,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, model.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return attempt, nil
}

func (r *AttemptRepository) Update(ctx context.Context, attempt *model.Attempt) error {
	query := `
		UPDATE attempts
		SET status = $1, completed_at = $2
		WHERE id = $3
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		attempt.Status,
		attempt.CompletedAt,
		attempt.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return model.ErrNotFound
	}

	return nil
}
