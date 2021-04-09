package hello

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/hello", Hello)

	return router
}

func Hello(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, "Hello, welcome")
}
