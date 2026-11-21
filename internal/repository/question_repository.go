package repository

import (
	"context"
	"database/sql"
	"errors"

	"qubix/internal/model"
)

type questionRepository struct {
	db *sql.DB
}

func NewQuestionRepository(db *sql.DB) *questionRepository {
	return &questionRepository{db: db}
}

func (r *questionRepository) Create(ctx context.Context, question *model.Question) error {
	query := `
		INSERT INTO questions (quiz_id, text, question_type, points)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	return r.db.QueryRowContext(
		ctx,
		query,
		question.QuizID,
		question.Text,
		question.QuestionType,
		question.Points,
	).Scan(&question.ID, &question.CreatedAt)
}

func (r *questionRepository) GetByID(ctx context.Context, id int) (*model.Question, error) {
	question := &model.Question{}

	query := `
		SELECT id, quiz_id, text, question_type, points, created_at
		FROM questions
		WHERE id = $1
	`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&question.ID,
		&question.QuizID,
		&question.Text,
		&question.QuestionType,
		&question.Points,
		&question.CreatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, model.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return question, nil
}

func (r *questionRepository) GetByQuizID(ctx context.Context, quizID int) ([]model.Question, error) {
	query := `
		SELECT id, quiz_id, text, question_type, points, created_at
		FROM questions
		WHERE quiz_id = $1
		ORDER BY id
	`

	rows, err := r.db.QueryContext(ctx, query, quizID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []model.Question

	for rows.Next() {
		var question model.Question

		err := rows.Scan(
			&question.ID,
			&question.QuizID,
			&question.Text,
			&question.QuestionType,
			&question.Points,
			&question.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		questions = append(questions, question)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return questions, nil
}

func (r *questionRepository) Update(ctx context.Context, question *model.Question) error {
	query := `
		UPDATE questions
		SET quiz_id = $1, text = $2, question_type = $3, points = $4
		WHERE id = $5
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		question.QuizID,
		question.Text,
		question.QuestionType,
		question.Points,
		question.ID,
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

func (r *questionRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM questions WHERE id = $1`

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
