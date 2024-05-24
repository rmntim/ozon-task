package models

import (
	"time"
)

type Comment struct {
	ID              uint      `json:"id"`
	Content         string    `json:"content"`
	AuthorID        uint      `json:"-" db:"author_id"`
	CreatedAt       time.Time `json:"createdAt" db:"created_at"`
	PostID          uint      `json:"-" db:"post_id"`
	ParentCommentID *uint     `json:"-" db:"parent_comment_id"`
	RepliesIDs      IDArray   `json:"-" db:"replies_ids"`
}

type Mutation struct {
}

type Post struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	Content     string    `json:"content"`
	AuthorID    uint      `json:"-" db:"author_id"`
	CommentsIDs IDArray   `json:"-" db:"comments_ids"`
}

type Query struct {
}

type Subscription struct {
}

type User struct {
	ID       uint    `json:"id"`
	Username string  `json:"username"`
	Email    string  `json:"email"`
	PostsIDs IDArray `json:"-" db:"posts_ids"`
}
