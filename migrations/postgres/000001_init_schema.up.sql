CREATE TABLE IF NOT EXISTS users
(
    id            SERIAL PRIMARY KEY,
    username      VARCHAR(50) UNIQUE NOT NULL,
    email         VARCHAR(50) UNIQUE NOT NULL,
    password_hash bytea              NOT NULL
);

CREATE TABLE IF NOT EXISTS posts
(
    id             SERIAL PRIMARY KEY,
    title          VARCHAR(50) UNIQUE NOT NULL,
    creator_id     INT                NOT NULL REFERENCES users,
    created_at     TIMESTAMPTZ        NOT NULL DEFAULT NOW(),
    content        TEXT               NOT NULL,
    allow_comments BOOLEAN            NOT NULL
);

CREATE TABLE IF NOT EXISTS comments
(
    id                SERIAL PRIMARY KEY,
    author_id         INT         NOT NULL REFERENCES users,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    likes             INT         NOT NULL DEFAULT 0,
    content           TEXT        NOT NULL,
    post_id           INT                  DEFAULT NULL,
    parent_comment_id INT                  DEFAULT NULL,
    FOREIGN KEY (post_id) REFERENCES posts ON DELETE CASCADE,
    FOREIGN KEY (parent_comment_id) REFERENCES comments ON DELETE CASCADE,
    CHECK (
        (post_id IS NOT NULL AND parent_comment_id IS NULL) OR
        (post_id IS NULL AND parent_comment_id IS NOT NULL)
        )
);