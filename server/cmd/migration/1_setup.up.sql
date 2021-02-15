DROP TABLE IF EXISTS chat_message;
CREATE TABLE chat_message (
	conversation_id TEXT,
	from_user TEXT,
	to_user TEXT,
	content TEXT,
	created_at TIMESTAMP
);

DROP TABLE IF EXISTS conversation;
CREATE TABLE conversation (
	id TEXT PRIMARY KEY,
	created_at TIMESTAMP
);

DROP TABLE IF EXISTS user_conversation;
CREATE TABLE user_conversation (
	conversation_id TEXT,
	user_name TEXT
);

DROP TABLE IF EXISTS chat_user;
CREATE TABLE chat_user
(
    id SERIAL NOT NULL,
    name TEXT NOT NULL UNIQUE,
    created_at timestamp DEFAULT now(),
);