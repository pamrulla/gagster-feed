package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/pamrulla/gagster-feed/pkg/hello"
)

func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(
		render.SetContentType(render.ContentTypeJSON), // Set Content-Type header as application/json
		middleware.Logger,          // Log API request calls
		middleware.RedirectSlashes, // Redirect slashes to no slash URL versions
		middleware.Recoverer,       // Recover from panics without crashing the server
	)

	router.Route("/v1", func(r chi.Router) {
		r.Mount("/api/hello", hello.Routes())
	})
	return router
}

func main() {
	router := Routes()

	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Printf("%s %s\n", method, route) // Walk and prints out all routes
		return nil
	}

	if err := chi.Walk(router, walkFunc); err != nil {
		log.Panicf("Logging err: %s\n", err.Error())
	}
	log.Fatal(http.ListenAndServe(":3000", router))
}
