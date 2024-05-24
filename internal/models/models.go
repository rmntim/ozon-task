package models

import (
	"time"
)

type Comment struct {
	ID              uint      `json:"id"`
	Content         string    `json:"content"`
	AuthorID        uint      `json:"-"`
	CreatedAt       time.Time `json:"createdAt"`
	PostID          uint      `json:"-"`
	ParentCommentID *uint     `json:"-"`
	RepliesIDs      []uint    `json:"-"`
}

type Mutation struct {
}

type Post struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	CreatedAt   time.Time `json:"createdAt"`
	Content     string    `json:"content"`
	Author      uint      `json:"-"`
	CommentsIDs []uint    `json:"-"`
}

type Query struct {
}

type Subscription struct {
}

type User struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	PostsIDs []uint `json:"-"`
}
