package repository

import (
	"context"
	"database/sql"
	"errors"

	"qubix/internal/model"
)

type optionRepository struct {
	db *sql.DB
}

func NewOptionRepository(db *sql.DB) *optionRepository {
	return &optionRepository{db: db}
}

func (r *optionRepository) Create(ctx context.Context, option *model.Option) error {
	query := `
		INSERT INTO options (question_id, text, is_correct)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	return r.db.QueryRowContext(
		ctx,
		query,
		option.QuestionID,
		option.Text,
		option.IsCorrect,
	).Scan(&option.ID)
}

func (r *optionRepository) GetByID(ctx context.Context, id int) (*model.Option, error) {
	option := &model.Option{}

	query := `
		SELECT id, question_id, text, is_correct
		FROM options
		WHERE id = $1
	`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&option.ID,
		&option.QuestionID,
		&option.Text,
		&option.IsCorrect,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, model.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return option, nil
}

func (r *optionRepository) GetByQuestionID(ctx context.Context, questionID int) ([]model.Option, error) {
	query := `
		SELECT id, question_id, text, is_correct
		FROM options
		WHERE question_id = $1
		ORDER BY id
	`

	rows, err := r.db.QueryContext(ctx, query, questionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var options []model.Option

	for rows.Next() {
		var option model.Option

		err := rows.Scan(
			&option.ID,
			&option.QuestionID,
			&option.Text,
			&option.IsCorrect,
		)
		if err != nil {
			return nil, err
		}

		options = append(options, option)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return options, nil
}

func (r *optionRepository) Update(ctx context.Context, option *model.Option) error {
	query := `
		UPDATE options
		SET question_id = $1, text = $2, is_correct = $3
		WHERE id = $4
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		option.QuestionID,
		option.Text,
		option.IsCorrect,
		option.ID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return model.ErrNotFound
	}

	return nil
}

func (r *optionRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM options WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return model.ErrNotFound
	}

	return nil
}
