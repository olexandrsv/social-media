CREATE TABLE projects (
	id SERIAL PRIMARY KEY,
	title VARCHAR(50),
	description TEXT,
	stack TEXT,
	user_id INTEGER
);

