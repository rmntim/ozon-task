CREATE TABLE IF NOT EXISTS users
(
    id            SERIAL PRIMARY KEY,
    username      VARCHAR(50)  NOT NULL UNIQUE,
    email         VARCHAR(100) NOT NULL UNIQUE,
    password_hash bytea        NOT NULL
);

CREATE TABLE IF NOT EXISTS posts
(
    id         SERIAL PRIMARY KEY,
    title      VARCHAR(255) NOT NULL,
    content    TEXT         NOT NULL,
    created_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    author_id  INTEGER      NOT NULL,
    FOREIGN KEY (author_id) REFERENCES users
);

CREATE TABLE IF NOT EXISTS comments
(
    id                SERIAL PRIMARY KEY,
    content           TEXT      NOT NULL,
    author_id         INTEGER   NOT NULL,
    created_at        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    post_id           INTEGER   NOT NULL,
    parent_comment_id INTEGER,
    FOREIGN KEY (author_id) REFERENCES users,
    FOREIGN KEY (post_id) REFERENCES posts,
    FOREIGN KEY (parent_comment_id) REFERENCES comments,
    CONSTRAINT parent_comment_in_same_post CHECK (
        parent_comment_id IS NULL
            OR post_id = (SELECT post_id
                          FROM comments
                          WHERE id = parent_comment_id)
        )
);

CREATE INDEX idx_author_id ON posts (author_id);
CREATE INDEX idx_post_id ON comments (post_id);
CREATE INDEX idx_parent_comment_id ON comments (parent_comment_id);