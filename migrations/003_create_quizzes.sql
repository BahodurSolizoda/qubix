CREATE TABLE quizzes (
    id SERIAL PRIMARY KEY,
    creator_id INTEGER NOT NULL REFERENCES users(id),
    category_id INTEGER NOT NULL REFERENCES categories(id),
    title TEXT NOT NULL,
    description TEXT,
    difficulty TEXT NOT NULL
        CHECK (difficulty IN ('easy', 'medium', 'hard')),
    time_limit INTEGER NOT NULL DEFAULT 0
        CHECK (time_limit >= 0),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);


