// Middleware for logging, authentication, etc.
package api

import (
	"log"
	"net/http"
	"time"
)

// LoggingMiddleware logs the HTTP method, path, status, and duration for each request.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(lrw, r)
		duration := time.Since(start)
		log.Printf("%s %s %d %s", r.Method, r.URL.Path, lrw.statusCode, duration)
	})
}

// loggingResponseWriter wraps http.ResponseWriter to capture status code.
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// TODOs for recommended middleware:
//
// 1. RequestIDMiddleware: Assigns a unique request ID to each request, adds it to context and response headers for traceability.
// 2. AuthMiddleware: Validates authentication tokens (e.g., JWT, API key) and rejects unauthorized requests.
// 3. RecoveryMiddleware: Catches panics, logs the error, and returns a 500 error response without crashing the server.
// 4. CORSMiddleware: Adds CORS headers to allow cross-origin requests from browsers.
// 5. RateLimitMiddleware: Limits the number of requests per client/IP to prevent abuse.
// 6. MetricsMiddleware: Collects and exposes metrics (request count, duration, errors) for monitoring.
// 7. ContentTypeMiddleware: Ensures requests have the correct Content-Type (e.g., application/json).
//
// Implement these as needed for production readiness and security.
