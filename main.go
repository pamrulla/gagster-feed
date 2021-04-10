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
		r.Route("/v1", func(r chi.Router) {
			r.Route("/users", func(r chi.Router) {
				r.Get("/", v1.GetUsers)
				r.Post("/", v1.CreateUser)
				r.Route("/{user_id}", func(r chi.Router) {
					r.Get("/", v1.GetUser)
					r.Put("/", v1.UpdateUser)
					r.Delete("/", v1.DeleteUser)
				})
				r.Put("/enable/{user_id}", v1.EnableUser)
				r.Put("/disable/{user_id}", v1.EnableUser)
			})
			r.Route("/gags/{user_id}", func(r chi.Router) {
				r.Get("/", v1.GetGags)
				r.Post("/", v1.CreateGag)
				r.Route("/{gag_id}", func(r chi.Router) {
					r.Get("/", v1.GetGag)
					r.Put("/", v1.UpdateGag)
					r.Delete("/", v1.DeleteGag)
				})
			})
			r.Put("/gags/enable/{gag_id}", v1.EnableUser)
			r.Put("/gags/disable/{gag_id}", v1.EnableUser)
			r.Route("/hearts/{gag_id}", func(r chi.Router) {
				r.Get("/", v1.GetHearts)
				r.Route("/{user_id}", func(r chi.Router) {
					r.Post("/", v1.CreateHeart)
					r.Delete("/", v1.DeleteHeart)
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
