package v1

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/pamrulla/gagster-feed/v1/gag"
	"github.com/pamrulla/gagster-feed/v1/heart"
	"github.com/pamrulla/gagster-feed/v1/user"
)

func Init() {
	user.Init()
	gag.Init()
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

// User Table Hanlders
func GetUsers(w http.ResponseWriter, r *http.Request) {
	user.GetUsers(w, r)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	user.Create(w, r)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	user.Update(w, r)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	user.Delete(w, r)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	user.Get(w, r)
}

func EnableUser(w http.ResponseWriter, r *http.Request) {
	user.Enable(w, r)
}

func DisableUser(w http.ResponseWriter, r *http.Request) {
	user.Disable(w, r)
}

// Gag Table Handler
func GetGags(w http.ResponseWriter, r *http.Request) {
	gag.GetGags(w, r)
}

func CreateGag(w http.ResponseWriter, r *http.Request) {
	gag.Create(w, r)
}

func UpdateGag(w http.ResponseWriter, r *http.Request) {
	gag.Update(w, r)
}

func DeleteGag(w http.ResponseWriter, r *http.Request) {
	gag.Delete(w, r)
}

func GetGag(w http.ResponseWriter, r *http.Request) {
	gag.Get(w, r)
}

func EnableGag(w http.ResponseWriter, r *http.Request) {
	gag.Enable(w, r)
}

func DisableGag(w http.ResponseWriter, r *http.Request) {
	gag.Disable(w, r)
}

// Hearts Table Handler
func GetHearts(w http.ResponseWriter, r *http.Request) {
	heart.GetHearts(w, r)
}

func CreateHeart(w http.ResponseWriter, r *http.Request) {
	heart.Create(w, r)
}

func DeleteHeart(w http.ResponseWriter, r *http.Request) {
	heart.Delete(w, r)
}
