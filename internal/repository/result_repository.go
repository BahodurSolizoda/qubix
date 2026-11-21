package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"qubix/internal/model"
)

type ResultRepository struct {
	db *sql.DB
}

func NewResultRepository(db *sql.DB) *ResultRepository {
	return &ResultRepository{
		db: db,
	}
}

func (r *ResultRepository) Create(ctx context.Context, result *model.Result) error {
	query := `
		INSERT INTO results (
			attempt_id,
			score,
			total_questions,
			correct_answers,
			is_passed
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`

	return r.db.QueryRowContext(
		ctx,
		query,
		result.AttemptID,
		result.Score,
		result.TotalQuestions,
		result.CorrectAnswers,
		result.IsPassed,
	).Scan(&result.ID, &result.CreatedAt)
}

func (r *ResultRepository) CreateAndCompleteAttempt(
	ctx context.Context,
	result *model.Result,
) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO results (
			attempt_id,
			score,
			total_questions,
			correct_answers,
			is_passed
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`

	err = tx.QueryRowContext(
		ctx,
		query,
		result.AttemptID,
		result.Score,
		result.TotalQuestions,
		result.CorrectAnswers,
		result.IsPassed,
	).Scan(&result.ID, &result.CreatedAt)
	if err != nil {
		return err
	}

	updateQuery := `
		UPDATE attempts
		SET status = 'completed', completed_at = $1
		WHERE id = $2 AND status = 'in_progress'
	`

	updateResult, err := tx.ExecContext(
		ctx,
		updateQuery,
		time.Now(),
		result.AttemptID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := updateResult.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return model.ErrAttemptCompleted
	}

	return tx.Commit()
}

func (r *ResultRepository) GetByAttemptID(
	ctx context.Context,
	attemptID int,
) (*model.Result, error) {
	result := &model.Result{}

	query := `
		SELECT id, attempt_id, score, total_questions,
		       correct_answers, is_passed, created_at
		FROM results
		WHERE attempt_id = $1
	`

	err := r.db.QueryRowContext(ctx, query, attemptID).Scan(
		&result.ID,
		&result.AttemptID,
		&result.Score,
		&result.TotalQuestions,
		&result.CorrectAnswers,
		&result.IsPassed,
		&result.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, model.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return result, nil
}
