package api

import (
	"log/slog"
	"net/http"
	"time"
)

// responseWriterWrapper is a wrapper for http.ResponseWriter that saves HTTP
// status code when it is set in the writer. And this code can be read after that
// request was sent to log it.
//
// If code was not set during request processing, it will 0.
type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (w *responseWriterWrapper) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *responseWriterWrapper) Write(body []byte) (int, error) {
	if w.statusCode == 0 {
		w.statusCode = http.StatusOK
	}

	return w.ResponseWriter.Write(body)
}

func loggingMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startingTime := time.Now()
		warpedWriter := &responseWriterWrapper{
			ResponseWriter: w,
		}

		handler(warpedWriter, r)
		logRESTRequest(warpedWriter, r, time.Since(startingTime))
	}
}

func logRESTRequest(w *responseWriterWrapper, r *http.Request, elapsedTime time.Duration) {
	lvl := levelByHTTPStatusCode(w.statusCode)
	slog.Log(r.Context(), lvl, "Request handled", slog.String("handler_type", "REST API"),
		slog.String("method", r.Method), slog.String("path", r.URL.Path), slog.Int("response_code", w.statusCode), slog.Duration("elapsed_time", elapsedTime))
}

func levelByHTTPStatusCode(code int) slog.Level {
	if code == http.StatusInternalServerError {
		return slog.LevelError
	}

	return slog.LevelInfo
}
