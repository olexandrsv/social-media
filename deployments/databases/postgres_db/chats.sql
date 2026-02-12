CREATE TABLE chats_users(
	chat_id INTEGER,
	user_id INTEGER,
	last_read_message TEXT
);

CREATE TABLE chats (
	name VARCHAR(50),
	id SERIAL PRIMARY KEY,
	owner INTEGER
);