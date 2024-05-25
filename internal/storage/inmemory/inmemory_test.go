package inmemory_test

import (
	"context"
	"github.com/rmntim/ozon-task/internal/storage/inmemory"
	"reflect"
	"testing"
)

func TestStorage_CreateComment(t *testing.T) {
	s := inmemory.New()

	ctx := context.Background()
	user, err := s.CreateUser(ctx, "test", "test", "test")
	if err != nil {
		t.Error("user should be created")
	}

	post, err := s.CreatePost(ctx, "test", "test", user.ID)
	if err != nil {
		t.Error("post should be created")
	}

	comment, err := s.CreateComment(ctx, "test", user.ID, post.ID, nil)
	if err != nil {
		t.Error("comment should be created")
	}

	if comment == nil {
		t.Error("comment should not be nil")
	}

	if comment.ID != 0 {
		t.Error("comment id should be 0")
	}

	if comment.Content != "test" {
		t.Error("comment content should be 'test'")
	}

	if comment.AuthorID != 1 {
		t.Error("comment author id should be 1")
	}

	if comment.PostID != 1 {
		t.Error("comment post id should be 1")
	}

	if comment.ParentCommentID != nil {
		t.Error("comment parent comment id should be nil")
	}
}

func TestStorage_GetCommentById(t *testing.T) {
	s := inmemory.New()

	ctx := context.Background()
	user, err := s.CreateUser(ctx, "test", "test", "test")
	if err != nil {
		t.Error("user should be created")
	}

	post, err := s.CreatePost(ctx, "test", "test", user.ID)
	if err != nil {
		t.Error("post should be created")
	}

	comment, err := s.CreateComment(ctx, "test", user.ID, post.ID, nil)
	if err != nil {
		t.Error("comment should be created")
	}

	if comment == nil {
		t.Error("comment should not be nil")
	}

	_, err = s.GetCommentById(ctx, comment.ID)
	if err != nil {
		t.Error("comment should be found")
	}
}

func TestStorage_CreatePost(t *testing.T) {
	s := inmemory.New()

	ctx := context.Background()
	user, err := s.CreateUser(ctx, "test", "test", "test")
	if err != nil {
		t.Error("user should be created")
	}

	post, err := s.CreatePost(ctx, "test", "test", user.ID)
	if err != nil {
		t.Error("post should be created")
	}

	if post == nil {
		t.Error("post should not be nil")
	}

	if post.ID != 0 {
		t.Error("post id should be 0")
	}

	if post.Title != "test" {
		t.Error("post title should be 'test'")
	}

	if post.Content != "test" {
		t.Error("post content should be 'test'")
	}

	if post.AuthorID != 1 {
		t.Error("post author id should be 1")
	}
}

func TestStorage_GetPostById(t *testing.T) {
	s := inmemory.New()

	ctx := context.Background()
	user, err := s.CreateUser(ctx, "test", "test", "test")
	if err != nil {
		t.Error("user should be created")
	}

	post, err := s.CreatePost(ctx, "test", "test", user.ID)
	if err != nil {
		t.Error("post should be created")
	}

	if post == nil {
		t.Error("post should not be nil")
	}

	_, err = s.GetPostById(ctx, post.ID)
	if err != nil {
		t.Error("post should be found")
	}
}

func TestStorage_CreateUser(t *testing.T) {
	s := inmemory.New()

	ctx := context.Background()
	user, err := s.CreateUser(ctx, "test", "test", "test")
	if err != nil {
		t.Error("user should be created")
	}

	if user == nil {
		t.Error("user should not be nil")
	}

	if user.ID != 0 {
		t.Error("user id should be 0")
	}

	if user.Username != "test" {
		t.Error("user username should be 'test'")
	}

	if user.Email != "test" {
		t.Error("user email should be 'test'")
	}
}

func TestStorage_GetUserById(t *testing.T) {
	s := inmemory.New()

	ctx := context.Background()
	user, err := s.CreateUser(ctx, "test", "test", "test")
	if err != nil {
		t.Error("user should be created")
	}

	if user == nil {
		t.Error("user should not be nil")
	}

	_, err = s.GetUserById(ctx, user.ID)
	if err != nil {
		t.Error("user should be found")
	}
}

func TestStorage_GetComments(t *testing.T) {
	s := inmemory.New()

	ctx := context.Background()
	user, err := s.CreateUser(ctx, "test", "test", "test")
	if err != nil {
		t.Error("user should be created")
	}

	post, err := s.CreatePost(ctx, "test", "test", user.ID)
	if err != nil {
		t.Error("post should be created")
	}

	comment, err := s.CreateComment(ctx, "test", user.ID, post.ID, nil)
	if err != nil {
		t.Error("comment should be created")
	}

	comments, err := s.GetComments(ctx, 10, 0)
	if err != nil {
		t.Error("comments should be found")
	}

	if len(comments) == 0 {
		t.Error("comments should not be empty")
	}

	if !reflect.DeepEqual(comments[0], comment) {
		t.Error("comments should be equal")
	}
}

func TestStorage_GetCommentsForPost(t *testing.T) {
	s := inmemory.New()

	ctx := context.Background()
	user, err := s.CreateUser(ctx, "test", "test", "test")
	if err != nil {
		t.Error("user should be created")
	}

	post, err := s.CreatePost(ctx, "test", "test", user.ID)
	if err != nil {
		t.Error("post should be created")
	}

	comment, err := s.CreateComment(ctx, "test", user.ID, post.ID, nil)
	if err != nil {
		t.Error("comment should be created")
	}

	comments, err := s.GetCommentsForPost(ctx, post.ID)
	if err != nil {
		t.Error("comments should be found")
	}

	if len(comments) == 0 {
		t.Error("comments should not be empty")
	}

	if !reflect.DeepEqual(comments[0], comment) {
		t.Error("comments should be equal")
	}
}

func TestStorage_GetPosts(t *testing.T) {
	s := inmemory.New()

	ctx := context.Background()
	user, err := s.CreateUser(ctx, "test", "test", "test")
	if err != nil {
		t.Error("user should be created")
	}

	post, err := s.CreatePost(ctx, "test", "test", user.ID)
	if err != nil {
		t.Error("post should be created")
	}

	posts, err := s.GetPosts(ctx, 10, 0)
	if err != nil {
		t.Error("posts should be found")
	}

	if len(posts) == 0 {
		t.Error("posts should not be empty")
	}

	if !reflect.DeepEqual(posts[0], post) {
		t.Error("posts should be equal")
	}
}

func TestStorage_GetPostsFromUser(t *testing.T) {
	s := inmemory.New()

	ctx := context.Background()
	user, err := s.CreateUser(ctx, "test", "test", "test")
	if err != nil {
		t.Error("user should be created")
	}

	post, err := s.CreatePost(ctx, "test", "test", 1)
	if err != nil {
		t.Error("post should be created")
	}

	posts, err := s.GetPostsFromUser(ctx, user.ID)
	if err != nil {
		t.Error("posts should be found")
	}

	if len(posts) == 0 {
		t.Error("posts should not be empty")
	}

	if !reflect.DeepEqual(posts[0], post) {
		t.Error("posts should be equal")
	}
}

func TestStorage_GetReplies(t *testing.T) {
	s := inmemory.New()

	ctx := context.Background()
	user, err := s.CreateUser(ctx, "test", "test", "test")
	if err != nil {
		t.Error("user should be created")
	}

	post, err := s.CreatePost(ctx, "test", "test", user.ID)
	if err != nil {
		t.Error("post should be created")
	}

	comment, err := s.CreateComment(ctx, "test", user.ID, post.ID, nil)
	if err != nil {
		t.Error("comment should be created")
	}

	replies, err := s.GetReplies(ctx, post.ID)
	if err != nil {
		t.Error("replies should be found")
	}

	if len(replies) == 0 {
		t.Error("replies should not be empty")
	}

	if !reflect.DeepEqual(replies[0], comment) {
		t.Error("replies should be equal")
	}
}

func TestStorage_GetUsers(t *testing.T) {
	s := inmemory.New()

	ctx := context.Background()
	user, err := s.CreateUser(ctx, "test", "test", "test")
	if err != nil {
		t.Error("user should be created")
	}

	users, err := s.GetUsers(ctx, 10, 0)
	if err != nil {
		t.Error("users should be found")
	}

	if len(users) == 0 {
		t.Error("users should not be empty")
	}

	if !reflect.DeepEqual(users[0], user) {
		t.Error("users should be equal")
	}
}

func TestStorage_ToggleComments(t *testing.T) {
	s := inmemory.New()

	ctx := context.Background()
	user, err := s.CreateUser(ctx, "test", "test", "test")
	if err != nil {
		t.Error("user should be created")
	}

	post, err := s.CreatePost(ctx, "test", "test", user.ID)
	if err != nil {
		t.Error("post should be created")
	}

	canComment, err := s.ToggleComments(ctx, post.ID, user.ID)
	if err != nil {
		t.Error("should not return error")
	}

	if canComment {
		t.Error("user should not be able to comment")
	}
}

func TestStorage_ToggleCommentsUnauthorized(t *testing.T) {
	s := inmemory.New()

	ctx := context.Background()
	user, err := s.CreateUser(ctx, "test", "test", "test")
	if err != nil {
		t.Error("user should be created")
	}

	post, err := s.CreatePost(ctx, "test", "test", user.ID)
	if err != nil {
		t.Error("post should be created")
	}

	_, err = s.ToggleComments(ctx, post.ID, 2)
	if err == nil {
		t.Error("should return error")
	}
}
