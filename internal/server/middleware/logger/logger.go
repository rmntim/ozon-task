package logger

import (
	"github.com/rmntim/ozon-task/internal/lib/logger/sl"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"time"
)

func New(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log.Info("Logger middleware is enabled")

		fn := func(w http.ResponseWriter, r *http.Request) {
			entry := log.With(
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
			)

			ww := httptest.NewRecorder()

			t1 := time.Now()
			defer func() {
				entry.Info("request completed",
					slog.Int("status", ww.Code),
					slog.Int("bytes", ww.Body.Len()),
					slog.String("duration", time.Since(t1).String()),
				)

				// copy everything from response recorder
				// to actual response writer
				for k, v := range ww.Result().Header {
					w.Header()[k] = v
				}
				w.WriteHeader(ww.Code)
				if _, err := ww.Body.WriteTo(w); err != nil {
					log.Error("Failed to write response body", sl.Err(err))
				}
			}()

			next.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(fn)
	}
}
