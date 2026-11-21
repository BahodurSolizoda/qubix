INSERT INTO categories (name, description)
VALUES
    ('Programming', 'Programming languages and development'),
    ('English', 'English language tests'),
    ('Mathematics', 'Mathematics tests'),
    ('History', 'History tests');

SELECT * FROM categories;

SELECT * FROM categories
ORDER BY id;