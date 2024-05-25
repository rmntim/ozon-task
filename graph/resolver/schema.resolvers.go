package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.47

import (
	"context"
	"errors"
	"github.com/rmntim/ozon-task/internal/lib/auth"
	"log/slog"

	"github.com/rmntim/ozon-task/graph"
	"github.com/rmntim/ozon-task/internal/lib/logger/sl"
	"github.com/rmntim/ozon-task/internal/lib/random"
	"github.com/rmntim/ozon-task/internal/models"
	"github.com/rmntim/ozon-task/internal/server"
)

var (
	postCreatedChannels  = make(map[string]chan *models.Post)
	commentAddedChannels = make(map[string]commentWithPost)
)

type commentWithPost struct {
	postID      uint
	commentChan chan *models.Comment
}

// Author is the resolver for the author field.
func (r *commentResolver) Author(ctx context.Context, obj *models.Comment) (*models.User, error) {
	const op = "resolver.Author"
	user, err := r.db.GetUserById(ctx, obj.AuthorID)
	if err != nil {
		r.log.Error("internal error", slog.String("op", op), sl.Err(err))
		return nil, server.ErrInternal
	}
	return user, nil
}

// Post is the resolver for the post field.
func (r *commentResolver) Post(ctx context.Context, obj *models.Comment) (*models.Post, error) {
	const op = "resolver.Post"
	post, err := r.db.GetPostById(ctx, obj.PostID)
	if err != nil {
		r.log.Error("internal error", slog.String("op", op), sl.Err(err))
		return nil, server.ErrInternal
	}
	return post, nil
}

// ParentComment is the resolver for the parentComment field.
func (r *commentResolver) ParentComment(ctx context.Context, obj *models.Comment) (*models.Comment, error) {
	const op = "resolver.ParentComment"
	if obj.ParentCommentID == nil {
		return nil, nil
	}
	parentComment, err := r.db.GetCommentById(ctx, *obj.ParentCommentID)
	if err != nil {
		r.log.Error("internal error", slog.String("op", op), sl.Err(err))
		return nil, server.ErrInternal
	}
	return parentComment, nil
}

// Replies is the resolver for the replies field.
func (r *commentResolver) Replies(ctx context.Context, obj *models.Comment) ([]*models.Comment, error) {
	const op = "resolver.Replies"
	replies, err := r.db.GetReplies(ctx, obj.ID)
	if err != nil {
		r.log.Error("internal error", slog.String("op", op), sl.Err(err))
		return nil, server.ErrInternal
	}
	return replies, nil
}

// CreateUser is the resolver for the createUser field.
func (r *mutationResolver) CreateUser(ctx context.Context, username string, email string, password string) (*models.User, error) {
	const op = "resolver.CreateUser"
	newUser, err := r.db.CreateUser(ctx, username, email, password)
	if err != nil {
		r.log.Error("internal error", slog.String("op", op), sl.Err(err))
		return nil, server.ErrInternal
	}
	return newUser, nil
}

// CreatePost is the resolver for the createPost field.
func (r *mutationResolver) CreatePost(ctx context.Context, title string, content string, authorID uint) (*models.Post, error) {
	const op = "resolver.CreatePost"
	newPost, err := r.db.CreatePost(ctx, title, content, authorID)
	if err != nil {
		r.log.Error("internal error", slog.String("op", op), sl.Err(err))
		return nil, server.ErrInternal
	}

	for _, observer := range postCreatedChannels {
		observer <- newPost
	}

	return newPost, nil
}

// CreateComment is the resolver for the createComment field.
func (r *mutationResolver) CreateComment(ctx context.Context, content string, authorID uint, postID uint, parentCommentID *uint) (*models.Comment, error) {
	const op = "resolver.CreateComment"
	newComment, err := r.db.CreateComment(ctx, content, authorID, postID, parentCommentID)
	if err != nil {
		if errors.Is(err, server.ErrCommentsDisabled) {
			return nil, server.ErrCommentsDisabled
		}
		r.log.Error("internal error", slog.String("op", op), sl.Err(err))
		return nil, server.ErrInternal
	}

	for _, observer := range commentAddedChannels {
		if postID == observer.postID {
			observer.commentChan <- newComment
		}
	}

	return newComment, nil
}

// ToggleComments is the resolver for the toggleComments field.
func (r *mutationResolver) ToggleComments(ctx context.Context, postID uint) (bool, error) {
	const op = "resolver.ToggleComments"
	user := auth.ForContext(ctx)
	if user == nil {
		return false, server.ErrUnauthorized
	}
	isEnabled, err := r.db.ToggleComments(ctx, postID, user.ID)
	if err != nil {
		r.log.Error("internal error", slog.String("op", op), sl.Err(err))
		return false, server.ErrInternal
	}
	return isEnabled, nil
}

// Author is the resolver for the author field.
func (r *postResolver) Author(ctx context.Context, obj *models.Post) (*models.User, error) {
	const op = "resolver.Author"
	user, err := r.db.GetUserById(ctx, obj.AuthorID)
	if err != nil {
		r.log.Error("internal error", slog.String("op", op), sl.Err(err))
		return nil, server.ErrInternal
	}
	return user, nil
}

// Comments is the resolver for the comments field.
func (r *postResolver) Comments(ctx context.Context, obj *models.Post) ([]*models.Comment, error) {
	const op = "resolver.Comments"
	comments, err := r.db.GetCommentsForPost(ctx, obj.ID)
	if err != nil {
		r.log.Error("internal error", slog.String("op", op), sl.Err(err))
		return nil, server.ErrInternal
	}
	return comments, nil
}

// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context, id uint) (*models.User, error) {
	const op = "resolver.User"
	user, err := r.db.GetUserById(ctx, id)
	if err != nil {
		if errors.Is(err, server.ErrUserNotFound) {
			return nil, server.ErrUserNotFound
		}
		r.log.Error("internal error", slog.String("op", op), sl.Err(err))
		return nil, server.ErrInternal
	}
	return user, nil
}

// Users is the resolver for the users field.
func (r *queryResolver) Users(ctx context.Context, limit int, offset int) ([]*models.User, error) {
	const op = "resolver.Users"
	users, err := r.db.GetUsers(ctx, limit, offset)
	if err != nil {
		r.log.Error("internal error", slog.String("op", op), sl.Err(err))
		return nil, server.ErrInternal
	}
	return users, nil
}

// Post is the resolver for the post field.
func (r *queryResolver) Post(ctx context.Context, id uint) (*models.Post, error) {
	const op = "resolver.Post"
	post, err := r.db.GetPostById(ctx, id)
	if err != nil {
		if errors.Is(err, server.ErrPostNotFound) {
			return nil, server.ErrPostNotFound
		}
		r.log.Error("internal error", slog.String("op", op), sl.Err(err))
		return nil, server.ErrInternal
	}
	return post, nil
}

// Posts is the resolver for the posts field.
func (r *queryResolver) Posts(ctx context.Context, limit int, offset int) ([]*models.Post, error) {
	const op = "resolver.Posts"
	posts, err := r.db.GetPosts(ctx, limit, offset)
	if err != nil {
		r.log.Error("internal error", slog.String("op", op), sl.Err(err))
		return nil, server.ErrInternal
	}
	return posts, nil
}

// Comment is the resolver for the comment field.
func (r *queryResolver) Comment(ctx context.Context, id uint) (*models.Comment, error) {
	const op = "resolver.Comment"
	comment, err := r.db.GetCommentById(ctx, id)
	if err != nil {
		if errors.Is(err, server.ErrCommentNotFound) {
			return nil, server.ErrCommentNotFound
		}
		r.log.Error("internal error", slog.String("op", op), sl.Err(err))
		return nil, server.ErrInternal
	}
	return comment, nil
}

// Comments is the resolver for the comments field.
func (r *queryResolver) Comments(ctx context.Context, limit int, offset int) ([]*models.Comment, error) {
	const op = "resolver.Comments"
	comments, err := r.db.GetComments(ctx, limit, offset)
	if err != nil {
		r.log.Error("internal error", slog.String("op", op), sl.Err(err))
		return nil, server.ErrInternal
	}
	return comments, nil
}

// PostAdded is the resolver for the postAdded field.
func (r *subscriptionResolver) PostAdded(ctx context.Context) (<-chan *models.Post, error) {
	id := random.NewRandomString(8)

	postEvent := make(chan *models.Post, 1)
	go func() {
		<-ctx.Done()
		close(postEvent)
		delete(postCreatedChannels, id)
	}()
	postCreatedChannels[id] = postEvent
	return postEvent, nil
}

// CommentAdded is the resolver for the commentAdded field.
func (r *subscriptionResolver) CommentAdded(ctx context.Context, postID uint) (<-chan *models.Comment, error) {
	id := random.NewRandomString(8)

	commentEvent := make(chan *models.Comment, 1)
	go func() {
		<-ctx.Done()
		close(commentEvent)
		delete(commentAddedChannels, id)
	}()
	commentAddedChannels[id] = commentWithPost{
		postID:      postID,
		commentChan: commentEvent,
	}
	return commentEvent, nil
}

// Posts is the resolver for the posts field.
func (r *userResolver) Posts(ctx context.Context, obj *models.User) ([]*models.Post, error) {
	const op = "resolver.Posts"
	posts, err := r.db.GetPostsFromUser(ctx, obj.ID)
	if err != nil {
		r.log.Error("internal error", slog.String("op", op), sl.Err(err))
		return nil, server.ErrInternal
	}
	return posts, nil
}

// Comment returns graph.CommentResolver implementation.
func (r *Resolver) Comment() graph.CommentResolver { return &commentResolver{r} }

// Mutation returns graph.MutationResolver implementation.
func (r *Resolver) Mutation() graph.MutationResolver { return &mutationResolver{r} }

// Post returns graph.PostResolver implementation.
func (r *Resolver) Post() graph.PostResolver { return &postResolver{r} }

// Query returns graph.QueryResolver implementation.
func (r *Resolver) Query() graph.QueryResolver { return &queryResolver{r} }

// Subscription returns graph.SubscriptionResolver implementation.
func (r *Resolver) Subscription() graph.SubscriptionResolver { return &subscriptionResolver{r} }

// User returns graph.UserResolver implementation.
func (r *Resolver) User() graph.UserResolver { return &userResolver{r} }

type commentResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type postResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
type userResolver struct{ *Resolver }
