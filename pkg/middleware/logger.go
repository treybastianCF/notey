package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

type statusResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// highjacking the og method
func (w *statusResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func Logger(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		wrapedWriter := &statusResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK, // default is ok
		}

		next.ServeHTTP(wrapedWriter, r)
		durration := time.Since(startTime)

		slog.Info("request",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("durration", durration.String()),
			slog.Int("status", wrapedWriter.statusCode))
	})
}
