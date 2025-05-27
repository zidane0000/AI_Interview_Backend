// API route definitions and HTTP server setup
package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// SetupRouter initializes the HTTP routes for the API using chi
func SetupRouter() http.Handler {
	r := chi.NewRouter()

	r.Use(LoggingMiddleware)

	// Custom NotFound for trailing slash
	r.NotFound(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/interviews/" {
			http.Error(w, ErrMsgMissingInterviewID, ErrCodeBadRequest)
			return
		}
		if r.URL.Path == "/evaluation/" {
			http.Error(w, ErrMsgMissingEvaluationID, ErrCodeBadRequest)
			return
		}
		http.NotFound(w, r)
	}))

	r.Route("/interviews", func(r chi.Router) {
		r.Post("/", CreateInterviewHandler)
		r.Get("/", ListInterviewsHandler)
		r.Get("/{id}", GetInterviewHandler)
	})

	r.Route("/evaluation", func(r chi.Router) {
		r.Post("/", SubmitEvaluationHandler)
		r.Get("/{id}", GetEvaluationHandler)
	})

	return r
}
