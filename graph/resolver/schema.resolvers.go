package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.47

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/rmntim/ozon-task/graph"
	"github.com/rmntim/ozon-task/graph/model"
	"github.com/rmntim/ozon-task/internal/lib/logger/sl"
)

// CreateUser is the resolver for the createUser field.
func (r *mutationResolver) CreateUser(ctx context.Context, username string, email string, password string) (*model.User, error) {
	const op = "resolver.CreateUser"
	newUser, err := r.db.CreateUser(ctx, username, email, password)
	if err != nil {
		r.log.Error("internal error", slog.String("op", op), sl.Err(err))
		return nil, ErrInternal
	}
	return newUser, nil
}

// CreatePost is the resolver for the createPost field.
func (r *mutationResolver) CreatePost(ctx context.Context, title string, content string, authorID uint) (*model.Post, error) {
	const op = "resolver.CreatePost"
	newPost, err := r.db.CreatePost(ctx, title, content, authorID)
	if err != nil {
		r.log.Error("internal error", slog.String("op", op), sl.Err(err))
		return nil, ErrInternal
	}
	return newPost, nil
}

// CreateComment is the resolver for the createComment field.
func (r *mutationResolver) CreateComment(ctx context.Context, content string, authorID uint, postID uint, parentCommentID *uint) (*model.Comment, error) {
	const op = "resolver.CreateComment"
	newComment, err := r.db.CreateComment(ctx, content, authorID, postID, parentCommentID)
	if err != nil {
		r.log.Error("internal error", slog.String("op", op), sl.Err(err))
		return nil, ErrInternal
	}
	return newComment, nil
}

// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context, id uint) (*model.User, error) {
	const op = "resolver.User"
	user, err := r.db.GetUserById(ctx, id)
	if err != nil {
		// FIXME: dont use sql errors in resolvers, but can't create generic errors in storage, as it causes cyclic imports
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		r.log.Error("internal error", slog.String("op", op), sl.Err(err))
		return nil, ErrInternal
	}
	return user, nil
}

// Users is the resolver for the users field.
func (r *queryResolver) Users(ctx context.Context, limit int, offset int) ([]*model.User, error) {
	const op = "resolver.Users"
	users, err := r.db.GetUsers(ctx, limit, offset)
	if err != nil {
		r.log.Error("internal error", slog.String("op", op), sl.Err(err))
		return nil, ErrInternal
	}
	return users, nil
}

// Post is the resolver for the post field.
func (r *queryResolver) Post(ctx context.Context, id uint) (*model.Post, error) {
	const op = "resolver.Post"
	post, err := r.db.GetPostById(ctx, id)
	if err != nil {
		r.log.Error("internal error", slog.String("op", op), sl.Err(err))
		return nil, ErrInternal
	}
	return post, nil
}

// Posts is the resolver for the posts field.
func (r *queryResolver) Posts(ctx context.Context, limit int, offset int) ([]*model.Post, error) {
	const op = "resolver.Posts"
	posts, err := r.db.GetPosts(ctx, limit, offset)
	if err != nil {
		r.log.Error("internal error", slog.String("op", op), sl.Err(err))
		return nil, ErrInternal
	}
	return posts, nil
}

// Comment is the resolver for the comment field.
func (r *queryResolver) Comment(ctx context.Context, id uint) (*model.Comment, error) {
	const op = "resolver.Comment"
	comment, err := r.db.GetCommentById(ctx, id)
	if err != nil {
		r.log.Error("internal error", slog.String("op", op), sl.Err(err))
		return nil, ErrInternal
	}
	return comment, nil
}

// Comments is the resolver for the comments field.
func (r *queryResolver) Comments(ctx context.Context, limit int, offset int) ([]*model.Comment, error) {
	const op = "resolver.Comments"
	comments, err := r.db.GetComments(ctx, limit, offset)
	if err != nil {
		r.log.Error("internal error", slog.String("op", op), sl.Err(err))
		return nil, ErrInternal
	}
	return comments, nil
}

// PostAdded is the resolver for the postAdded field.
func (r *subscriptionResolver) PostAdded(ctx context.Context) (<-chan *model.Post, error) {
	panic(fmt.Errorf("not implemented: PostAdded - postAdded"))
}

// CommentAdded is the resolver for the commentAdded field.
func (r *subscriptionResolver) CommentAdded(ctx context.Context, postID *uint) (<-chan *model.Comment, error) {
	panic(fmt.Errorf("not implemented: CommentAdded - commentAdded"))
}

// Mutation returns graph.MutationResolver implementation.
func (r *Resolver) Mutation() graph.MutationResolver { return &mutationResolver{r} }

// Query returns graph.QueryResolver implementation.
func (r *Resolver) Query() graph.QueryResolver { return &queryResolver{r} }

// Subscription returns graph.SubscriptionResolver implementation.
func (r *Resolver) Subscription() graph.SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//   - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//     it when you're done.
//   - You have helper methods in this file. Move them out to keep these resolver files clean.
var (
	ErrInternal     = errors.New("internal server error")
	ErrUserNotFound = errors.New("no such user")
)
