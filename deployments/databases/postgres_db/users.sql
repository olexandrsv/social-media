CREATE TABLE users (
	login varchar(50) NOT NULL,
	first_name varchar(50) NOT NULL,
	second_name varchar(50) NOT NULL,
	password varchar(100),
	bio TEXT DEFAULT '',
	interests TEXT  DEFAULT '',
	id SERIAL PRIMARY KEY
);

INSERT INTO users (login,first_name,second_name,password,bio,interests) VALUES 
	('bob','Bob','Smith','+gTsuWYzmHxGHoWFpxgwHYxj4ytP9WmEB1QVVx5EwaaxnWWvf987MaQXMSfr40kQ','student','c++, golang, java'),
	('ben','Ben','Jones','+gTsuWYzmHxGHoWFpxgwHYxj4ytP9WmEB1QVVx5EwaaxnWWvf987MaQXMSfr40kQ','pupil','java, html');


CREATE TABLE followers (
	last_read_post TEXT,
	user_id INTEGER,
	follower_id INTEGER
);