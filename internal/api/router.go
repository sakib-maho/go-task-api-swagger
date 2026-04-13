package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(handler *Handler) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", handler.Health)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/tasks", func(r chi.Router) {
			r.Get("/", handler.ListTasks)
			r.Post("/", handler.CreateTask)
			r.Get("/{taskID}", handler.GetTask)
			r.Put("/{taskID}", handler.UpdateTask)
			r.Delete("/{taskID}", handler.DeleteTask)
		})
	})

	r.Get("/swagger", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "docs/swagger.html")
	})
	r.Get("/swagger/openapi.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "docs/openapi.json")
	})

	return r
}
