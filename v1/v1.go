package v1

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/pamrulla/gagster-feed/v1/user"
)

func Init() {
	user.Init()
}

func Routes() *chi.Mux {
	//user.Init()
	router := chi.NewRouter()
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, "test")
	})
	//router.Get("/users", user.GetUsers)
	fmt.Printf("%+v\n", router.Routes())
	return router
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	user.GetUsers(w, r)
}
