CREATE TABLE IF NOT EXISTS users (
    uuid       VARCHAR(36)  PRIMARY KEY,
    email      VARCHAR(256) UNIQUE NOT NULL,
    passhash   VARCHAR(72)  NOT NULL,
    created_at TIMESTAMP    NOT NULL,
    updated_at TIMESTAMP    NOT NULL
);
