package auth

import (
	"context"
	"errors"
	"github.com/rmntim/ozon-task/internal/models"
	"github.com/rmntim/ozon-task/internal/server"
	"github.com/rmntim/ozon-task/internal/storage"
	"net/http"
	"strconv"
)

var userCtxKey = &contextKey{}

type contextKey struct{}

func Middleware(db storage.Storage) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c, err := r.Cookie("auth-cookie")

			if err != nil || c == nil {
				next.ServeHTTP(w, r)
				return
			}

			userId, err := validateAndGetUserID(c)
			if err != nil {
				http.Error(w, "Invalid cookie", http.StatusForbidden)
				return
			}

			user, err := db.GetUserById(r.Context(), uint(userId))
			if err != nil {
				if errors.Is(err, server.ErrUserNotFound) {
					http.Error(w, "no such user", http.StatusNotFound)
					return
				}
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}

			ctx := context.WithValue(r.Context(), userCtxKey, user)

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func validateAndGetUserID(c *http.Cookie) (int, error) {
	return strconv.Atoi(c.Value)
}

// ForContext finds the user from the context. REQUIRES Middleware to have run.
func ForContext(ctx context.Context) *models.User {
	raw, _ := ctx.Value(userCtxKey).(*models.User)
	return raw
}
