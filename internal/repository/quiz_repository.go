package repository

import (
	"context"
	"database/sql"
	"errors"

	"qubix/internal/model"
)

type QuizRepository struct {
	db *sql.DB
}

func NewQuizRepository(db *sql.DB) *QuizRepository {
	return &QuizRepository{
		db: db,
	}
}

func (r *QuizRepository) Create(ctx context.Context, quiz *model.Quiz) error {
	query := `
		INSERT INTO quizzes (
			creator_id, category_id, title, description,
			difficulty, time_limit
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRowContext(
		ctx,
		query,
		quiz.CreatorID,
		quiz.CategoryID,
		quiz.Title,
		quiz.Description,
		quiz.Difficulty,
		quiz.TimeLimit,
	).Scan(&quiz.ID, &quiz.CreatedAt, &quiz.UpdatedAt)
}

func (r *QuizRepository) GetByID(ctx context.Context, id int) (*model.Quiz, error) {
	quiz := &model.Quiz{}

	query := `
		SELECT id, creator_id, category_id, title, description,
		       difficulty, time_limit, created_at, updated_at
		FROM quizzes
		WHERE id = $1
	`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&quiz.ID,
		&quiz.CreatorID,
		&quiz.CategoryID,
		&quiz.Title,
		&quiz.Description,
		&quiz.Difficulty,
		&quiz.TimeLimit,
		&quiz.CreatedAt,
		&quiz.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, model.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return quiz, nil
}

func (r *QuizRepository) GetAll(ctx context.Context) ([]model.Quiz, error) {
	query := `
		SELECT id, creator_id, category_id, title, description,
		       difficulty, time_limit, created_at, updated_at
		FROM quizzes
		ORDER BY id
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var quizzes []model.Quiz

	for rows.Next() {
		var quiz model.Quiz

		err := rows.Scan(
			&quiz.ID,
			&quiz.CreatorID,
			&quiz.CategoryID,
			&quiz.Title,
			&quiz.Description,
			&quiz.Difficulty,
			&quiz.TimeLimit,
			&quiz.CreatedAt,
			&quiz.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		quizzes = append(quizzes, quiz)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return quizzes, nil
}

func (r *QuizRepository) Update(ctx context.Context, quiz *model.Quiz) error {
	query := `
		UPDATE quizzes
		SET category_id = $1,
		    title = $2,
		    description = $3,
		    difficulty = $4,
		    time_limit = $5,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $6
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		quiz.CategoryID,
		quiz.Title,
		quiz.Description,
		quiz.Difficulty,
		quiz.TimeLimit,
		quiz.ID,
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

func (r *QuizRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM quizzes WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
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
