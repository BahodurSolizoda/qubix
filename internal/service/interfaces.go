package service

import (
	"context"

	"qubix/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id int) (*model.User, error)
	GetByLoginOrEmail(ctx context.Context, value string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id int) error
}

type CategoryRepository interface {
	Create(ctx context.Context, category *model.Category) error
	GetByID(ctx context.Context, id int) (*model.Category, error)
	GetAll(ctx context.Context) ([]model.Category, error)
	Update(ctx context.Context, category *model.Category) error
	Delete(ctx context.Context, id int) error
}

type QuestionRepository interface {
	Create(ctx context.Context, question *model.Question) error
	GetByID(ctx context.Context, id int) (*model.Question, error)
	GetByQuizID(ctx context.Context, quizID int) ([]model.Question, error)
	Update(ctx context.Context, question *model.Question) error
	Delete(ctx context.Context, id int) error
}

type QuizRepository interface {
	Create(ctx context.Context, quiz *model.Quiz) error
	GetByID(ctx context.Context, id int) (*model.Quiz, error)
	GetAll(ctx context.Context) ([]model.Quiz, error)
	Update(ctx context.Context, quiz *model.Quiz) error
	Delete(ctx context.Context, id int) error
}

type OptionRepository interface {
	Create(ctx context.Context, option *model.Option) error
	GetByID(ctx context.Context, id int) (*model.Option, error)
	GetByQuestionID(ctx context.Context, questionID int) ([]model.Option, error)
	Update(ctx context.Context, option *model.Option) error
	Delete(ctx context.Context, id int) error
}

type AttemptRepository interface {
	Create(ctx context.Context, attempt *model.Attempt) error
	GetByID(ctx context.Context, id int) (*model.Attempt, error)
	Update(ctx context.Context, attempt *model.Attempt) error
}

type AnswerRepository interface {
	Create(ctx context.Context, answer *model.Answer) error
	CreateWithOptions(ctx context.Context, answer *model.Answer, optionIDs []int) error
	AddOption(ctx context.Context, answerOption *model.AnswerOption) error
	GetByAttemptID(ctx context.Context, attemptID int) ([]model.Answer, error)
	GetOptionIDsByAnswerID(ctx context.Context, answerID int) ([]int, error)
}

type ResultRepository interface {
	Create(ctx context.Context, result *model.Result) error
	CreateAndCompleteAttempt(ctx context.Context, result *model.Result) error
	GetByAttemptID(ctx context.Context, attemptID int) (*model.Result, error)
}
