CREATE TABLE IF NOT EXISTS users
(
    id            SERIAL PRIMARY KEY,
    username      VARCHAR(50)  NOT NULL UNIQUE,
    email         VARCHAR(100) NOT NULL UNIQUE,
    password_hash bytea        NOT NULL
);

CREATE TABLE IF NOT EXISTS posts
(
    id                 SERIAL PRIMARY KEY,
    title              VARCHAR(255) NOT NULL,
    content            TEXT         NOT NULL,
    created_at         TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    author_id          INTEGER      NOT NULL,
    comments_available BOOLEAN      NOT NULL DEFAULT TRUE,
    FOREIGN KEY (author_id) REFERENCES users
);

CREATE TABLE IF NOT EXISTS comments
(
    id                SERIAL PRIMARY KEY,
    content           VARCHAR(2000) NOT NULL,
    author_id         INTEGER       NOT NULL,
    created_at        TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    post_id           INTEGER       NOT NULL,
    parent_comment_id INTEGER,
    FOREIGN KEY (author_id) REFERENCES users,
    FOREIGN KEY (post_id) REFERENCES posts,
    FOREIGN KEY (parent_comment_id) REFERENCES comments
);

CREATE INDEX idx_author_id ON posts (author_id);
CREATE INDEX idx_post_id ON comments (post_id);
CREATE INDEX idx_parent_comment_id ON comments (parent_comment_id);

CREATE OR REPLACE FUNCTION check_post_for_comments() RETURNS trigger AS
$check_post_for_comments$
DECLARE
    post_comments_available BOOLEAN;
BEGIN
    SELECT comments_available INTO post_comments_available FROM posts WHERE id = NEW.post_id;

    IF NOT post_comments_available THEN
        RAISE EXCEPTION 'Post with ID % does not allow comments', NEW.post_id;
    END IF;

    RETURN NEW;
END;
$check_post_for_comments$ LANGUAGE plpgsql;

CREATE TRIGGER add_comment_with_check
    BEFORE INSERT
    ON comments
    FOR EACH ROW
EXECUTE PROCEDURE check_post_for_comments();

CREATE FUNCTION check_comment_same_post_as_parent() RETURNS trigger AS
$check_comment_same_post_as_parent$
DECLARE
    post_id INTEGER;
BEGIN
    IF NEW.parent_comment_id IS NOT NULL THEN
        SELECT post_id INTO post_id FROM comments WHERE id = NEW.parent_comment_id;

        IF NEW.post_id <> post_id THEN
            RAISE EXCEPTION 'Parent comment with ID % does not belong to post with ID %', NEW.parent_comment_id, NEW.post_id;
        END IF;
    END IF;
END;
$check_comment_same_post_as_parent$ LANGUAGE plpgsql;

CREATE TRIGGER comment_same_post_as_parent
    BEFORE INSERT
    ON comments
    FOR EACH ROW
EXECUTE PROCEDURE check_comment_same_post_as_parent();