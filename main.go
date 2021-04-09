package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	hello "github.com/pamrulla/gagster-feed/pkg"
	v1 "github.com/pamrulla/gagster-feed/v1"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func Routes() *chi.Mux {
	v1.Init()

	router := chi.NewRouter()
	router.Use(
		render.SetContentType(render.ContentTypeJSON), // Set Content-Type header as application/json
		middleware.Logger,          // Log API request calls
		middleware.RedirectSlashes, // Redirect slashes to no slash URL versions
		middleware.Recoverer,       // Recover from panics without crashing the server
	)

	router.Route("/api", func(r chi.Router) {
		r.Mount("/hello", hello.Routes())
		r.Route("/v1", func(r1 chi.Router) {
			r1.Route("/users", func(r chi.Router) {
				r.Get("/", v1.GetUsers)
				r.Route("/{user_id}", func(r chi.Router) {
					r.Get("/", v1.GetUser)
					r.Post("/", v1.CreateUser)
					r.Put("/", v1.UpdateUser)
					r.Delete("/", v1.DeleteUser)
				})
			})
		})
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
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}
