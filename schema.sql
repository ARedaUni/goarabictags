CREATE TABLE IF NOT EXISTS users (
    email VARCHAR(255) NOT NULL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    password_hash CHAR(60) NOT NULL,
    created TIMESTAMP NOT NULL,
    updated TIMESTAMP NOT NULL,

    CONSTRAINT user_username_uc UNIQUE (username)
);

CREATE TABLE IF NOT EXISTS excerpt (
    id SERIAL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    author_email VARCHAR(255) NOT NULL,
    created TIMESTAMP NOT NULL,
    updated TIMESTAMP NOT NULL,

    FOREIGN KEY (author_email)
        REFERENCES users(email)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS word (
    id SERIAL PRIMARY KEY,
    word VARCHAR(30) NOT NULL,
    word_pos INTEGER NOT NULL,
    connected BOOLEAN NOT NULL,
    punctuation BOOLEAN NOT NULL,
    excerpt_id INTEGER NOT NULL,
    created TIMESTAMP NOT NULL,
    updated TIMESTAMP NOT NULL,

    na_ignore BOOLEAN NOT NULL,
    na_sentence_start BOOLEAN NOT NULL,

    irab_case VARCHAR(30) NOT NULL,
    irab_state VARCHAR(30) NOT NULL,

    FOREIGN KEY (excerpt_id)
        REFERENCES excerpt(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS manuscript (
    id SERIAL PRIMARY KEY,
    content TEXT NOT NULL,
    excerpt_id INTEGER NOT NULL,
    created TIMESTAMP NOT NULL,
    updated TIMESTAMP NOT NULL,

    FOREIGN KEY (excerpt_id)
        REFERENCES excerpt(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT manuscript_excerpt_id_uc UNIQUE(excerpt_id)
);

CREATE TABLE IF NOT EXISTS sessions (
    token CHAR(43) PRIMARY KEY,
    data BYTEA NOT NULL,
    expiry TIMESTAMP NOT NULL
);

CREATE INDEX sessions_expiry_idx ON sessions (expiry);
