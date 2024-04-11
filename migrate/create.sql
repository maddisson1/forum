CREATE TABLE IF NOT EXISTS posts (
    post_id  INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    created DATETIME NOT NULL,
    user_id INTEGER,
    FOREIGN KEY (user_id) REFERENCES sessions (user_id)
);

CREATE TABLE IF NOT EXISTS users (
    user_id  INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE CONSTRAINT users_uc_username,
    email TEXT NOT NULL UNIQUE CONSTRAINT users_uc_email,
    hash_password TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS comments (
		comment_id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		post_id INTEGER NOT NULL,
		comment TEXT NOT NULL,
		created_at DATE NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(user_id),
		FOREIGN KEY (post_id) REFERENCES posts(post_id)
	);

CREATE TABLE IF NOT EXISTS sessions (
		user_id INTEGER,
		session_token TEXT,
		expires_at TIME
	);	

CREATE TABLE IF NOT EXISTS reactions (
		user_id INTEGER NOT NULL,
		post_id INTEGER NOT NULL,
		reaction_status INTEGER DEFAULT 0 NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(user_id),
		FOREIGN KEY (post_id) REFERENCES posts(post_id)
		PRIMARY KEY (user_id, post_id)
	);

CREATE TABLE IF NOT EXISTS comment_reactions (
		user_id INTEGER NOT NULL,
		comment_id INTEGER NOT NULL,
		reaction_status INTEGER DEFAULT 0 NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(user_id),
		FOREIGN KEY (comment_id) REFERENCES comments(comment_id)
		PRIMARY KEY (user_id, comment_id)
	);