CREATE TABLE answer_options (
    answer_id INTEGER NOT NULL REFERENCES answers(id),
    option_id INTEGER NOT NULL REFERENCES options(id),
    PRIMARY KEY (answer_id, option_id)
);


