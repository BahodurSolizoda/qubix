CREATE TABLE questions (
    id SERIAL PRIMARY KEY,
    quiz_id INTEGER NOT NULL REFERENCES quizzes(id),
    text TEXT NOT NULL,
    question_type TEXT NOT NULL
        CHECK (question_type IN ('single', 'multiple', 'text')),
    points INTEGER NOT NULL DEFAULT 1
        CHECK (points > 0),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);


