package server

import "errors"

var (
	ErrInternal        = errors.New("internal server error")
	ErrUserNotFound    = errors.New("no such user")
	ErrPostNotFound    = errors.New("no such post")
	ErrCommentNotFound = errors.New("no such comment")
)
