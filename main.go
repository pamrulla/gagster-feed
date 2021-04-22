package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/pamrulla/gagster-feed/helpers"
	hlp "github.com/pamrulla/gagster-feed/helpers"
	hello "github.com/pamrulla/gagster-feed/pkg"
	v1 "github.com/pamrulla/gagster-feed/v1"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
)

func Routes() *chi.Mux {
	v1.Init()
	helpers.InitiCloudinary()

	// router.Use(Cors)
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"X-PINGOTHER", "Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
		Debug:            true,
	})
	router := chi.NewRouter()
	router.Use(
		render.SetContentType(render.ContentTypeJSON), // Set Content-Type header as application/json
		middleware.Logger, // Log API request calls
		cors.Handler,
		middleware.RedirectSlashes, // Redirect slashes to no slash URL versions
		middleware.Recoverer,       // Recover from panics without crashing the server
	)

	router.Route("/api", func(r chi.Router) {
		r.Mount("/hello", hello.Routes())
		r.Route("/v1", func(r chi.Router) {
			r.Post("/login", v1.Login)
			r.Route("/users", func(r chi.Router) {
				r.Post("/", v1.CreateUser)
				r.Group(func(r chi.Router) {
					r.Use(jwtauth.Verifier(hlp.GetTokenAuth()))
					r.Use(jwtauth.Authenticator)
					r.Get("/", v1.GetUsers)
					r.Route("/{user_id}", func(r chi.Router) {
						r.Get("/", v1.GetUser)
						r.Put("/", v1.UpdateUser)
						r.Delete("/", v1.DeleteUser)
					})
					r.Put("/enable/{user_id}", v1.EnableUser)
					r.Put("/disable/{user_id}", v1.DisableUser)
				})
			})
			r.Get("/gags/feed", v1.Feed)
			r.Get("/gags/author/{user_id}", v1.GetAuthorGags)
			r.Get("/gags/{gag_id}", v1.GetGag)
			r.Get("/gags/tags", v1.GagsWithTags)
			r.Group(func(r chi.Router) {
				r.Use(jwtauth.Verifier(hlp.GetTokenAuth()))
				r.Use(jwtauth.Authenticator)
				r.Post("/gags/{user_id}", v1.CreateGag)
				r.Route("/gags", func(r chi.Router) {
					r.Route("/{gag_id}", func(r chi.Router) {
						r.Put("/", v1.UpdateGag)
						r.Delete("/", v1.DeleteGag)
					})
				})
				r.Put("/gags/enable/{gag_id}", v1.EnableUser)
				r.Put("/gags/disable/{gag_id}", v1.EnableUser)
			})

			r.Route("/hearts/{gag_id}", func(r chi.Router) {
				r.Get("/", v1.GetHearts)
				r.Group(func(r chi.Router) {
					r.Use(jwtauth.Verifier(hlp.GetTokenAuth()))
					r.Use(jwtauth.Authenticator)
					r.Route("/{user_id}", func(r chi.Router) {
						r.Post("/", v1.CreateHeart)
						r.Delete("/", v1.DeleteHeart)
					})
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
