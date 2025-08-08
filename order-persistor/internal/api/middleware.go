package api

import (
	"log/slog"
	"net/http"
	"time"
)

type Middleware func(http.Handler) http.Handler

func stackMiddleware(h http.Handler, mws ...Middleware) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}

	return h
}

func NewLogMiddleware(logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			wrappedWriter := &responseWriterWrapper{
				code:    200,
				wrapped: w,
			}

			next.ServeHTTP(wrappedWriter, r)

			duration := time.Since(start)

			fields := []any{
				"method", r.Method,
				"path", r.URL.Path,
				"status", wrappedWriter.code,
				"duration", duration.String(),
			}

			if wrappedWriter.code >= 400 {
				logger.Error("HTTP", fields...)
			} else {
				logger.Info("HTTP", fields...)
			}
		})
	}
}

var _ http.ResponseWriter = &responseWriterWrapper{}

type responseWriterWrapper struct {
	code    int
	wrapped http.ResponseWriter
}

func (l *responseWriterWrapper) Header() http.Header {
	return l.wrapped.Header()
}

func (l *responseWriterWrapper) Write(buf []byte) (int, error) {
	return l.wrapped.Write(buf)
}

func (l *responseWriterWrapper) WriteHeader(statusCode int) {
	l.code = statusCode
	l.wrapped.WriteHeader(statusCode)
}
