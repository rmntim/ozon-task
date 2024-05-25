package auth_test

import (
	"github.com/rmntim/ozon-task/internal/lib/auth"
	"github.com/rmntim/ozon-task/internal/storage/inmemory"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestForContext(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := auth.ForContext(r.Context())
		if user == nil {
			t.Error("user is nil")
		}

		if user.ID != 1 {
			t.Error("user id is not 1")
		}
	})

	testHandler := auth.Middleware(inmemory.New())(nextHandler)
	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{
		Name:  "auth-cookie",
		Value: "1",
	})
	testHandler.ServeHTTP(httptest.NewRecorder(), req)
}
