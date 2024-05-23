CREATE TABLE users
(
    id            SERIAL PRIMARY KEY,
    username      VARCHAR(50)  NOT NULL UNIQUE,
    email         VARCHAR(100) NOT NULL UNIQUE,
    password_hash bytea        NOT NULL
);

CREATE TABLE posts
(
    id        SERIAL PRIMARY KEY,
    title     VARCHAR(255) NOT NULL,
    content   TEXT         NOT NULL,
    author_id INTEGER      NOT NULL,
    FOREIGN KEY (author_id) REFERENCES users
);

CREATE TABLE comments
(
    id                SERIAL PRIMARY KEY,
    content           TEXT    NOT NULL,
    author_id         INTEGER NOT NULL,
    post_id           INTEGER NOT NULL,
    parent_comment_id INTEGER,
    FOREIGN KEY (author_id) REFERENCES users,
    FOREIGN KEY (post_id) REFERENCES posts,
    FOREIGN KEY (parent_comment_id) REFERENCES comments
);

CREATE INDEX idx_author_id ON posts (author_id);
CREATE INDEX idx_post_id ON comments (post_id);
CREATE INDEX idx_parent_comment_id ON comments (parent_comment_id);