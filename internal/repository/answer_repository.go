package repository

import (
	"context"
	"database/sql"

	"qubix/internal/model"
)

type AnswerRepository struct {
	db *sql.DB
}

func NewAnswerRepository(db *sql.DB) *AnswerRepository {
	return &AnswerRepository{db: db}
}

func (r *AnswerRepository) Create(ctx context.Context, answer *model.Answer) error {
	query := `
		INSERT INTO answers (attempt_id, question_id, text_answer)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	return r.db.QueryRowContext(
		ctx,
		query,
		answer.AttemptID,
		answer.QuestionID,
		answer.TextAnswer,
	).Scan(&answer.ID, &answer.CreatedAt)
}

func (r *AnswerRepository) CreateWithOptions(
	ctx context.Context,
	answer *model.Answer,
	optionIDs []int,
) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO answers (attempt_id, question_id, text_answer)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	err = tx.QueryRowContext(
		ctx,
		query,
		answer.AttemptID,
		answer.QuestionID,
		answer.TextAnswer,
	).Scan(&answer.ID, &answer.CreatedAt)
	if err != nil {
		return err
	}

	for _, optionID := range optionIDs {
		_, err = tx.ExecContext(
			ctx,
			`INSERT INTO answer_options (answer_id, option_id) VALUES ($1, $2)`,
			answer.ID,
			optionID,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *AnswerRepository) AddOption(
	ctx context.Context,
	answerOption *model.AnswerOption,
) error {
	query := `
		INSERT INTO answer_options (answer_id, option_id)
		VALUES ($1, $2)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		answerOption.AnswerID,
		answerOption.OptionID,
	)

	return err
}

func (r *AnswerRepository) GetByAttemptID(
	ctx context.Context,
	attemptID int,
) ([]model.Answer, error) {
	query := `
		SELECT id, attempt_id, question_id, text_answer, created_at
		FROM answers
		WHERE attempt_id = $1
		ORDER BY id
	`

	rows, err := r.db.QueryContext(ctx, query, attemptID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var answers []model.Answer

	for rows.Next() {
		var answer model.Answer

		err := rows.Scan(
			&answer.ID,
			&answer.AttemptID,
			&answer.QuestionID,
			&answer.TextAnswer,
			&answer.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		answers = append(answers, answer)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return answers, nil
}

func (r *AnswerRepository) GetOptionIDsByAnswerID(
	ctx context.Context,
	answerID int,
) ([]int, error) {
	query := `
		SELECT option_id
		FROM answer_options
		WHERE answer_id = $1
		ORDER BY option_id
	`

	rows, err := r.db.QueryContext(ctx, query, answerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var optionIDs []int

	for rows.Next() {
		var optionID int

		if err := rows.Scan(&optionID); err != nil {
			return nil, err
		}

		optionIDs = append(optionIDs, optionID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return optionIDs, nil
}
