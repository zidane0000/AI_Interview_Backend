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

// CORSMiddleware adds CORS headers to allow cross-origin requests from browsers
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// Development: Allow localhost origins
		// TODO: In production, replace with specific allowed origins
		allowedOrigins := []string{
			"http://localhost:3000",
			"http://localhost:5173",
			"http://127.0.0.1:3000",
			"http://127.0.0.1:5173",
		}

		// Check if origin is allowed
		isAllowed := false
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				isAllowed = true
				break
			}
		}

		if isAllowed {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		} else if origin == "" {
			// Allow same-origin requests (no Origin header)
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Type")
		w.Header().Set("Access-Control-Max-Age", "86400")

		// Handle preflight OPTIONS requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// TODOs for recommended middleware:
//
// 1. RequestIDMiddleware: Assigns a unique request ID to each request, adds it to context and response headers for traceability.
// 2. AuthMiddleware: Validates authentication tokens (e.g., JWT, API key) and rejects unauthorized requests.
// 3. RecoveryMiddleware: Catches panics, logs the error, and returns a 500 error response without crashing the server.
// 4. RateLimitMiddleware: Limits the number of requests per client/IP to prevent abuse.
// 5. MetricsMiddleware: Collects and exposes metrics (request count, duration, errors) for monitoring.
// 6. ContentTypeMiddleware: Ensures requests have the correct Content-Type (e.g., application/json).
//
// Implement these as needed for production readiness and security.
