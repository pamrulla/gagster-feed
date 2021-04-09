package user

import (
	"net/http"

	"github.com/go-chi/render"

	"github.com/pamrulla/gagster-feed/models"
)

var users models.Users

func Init() {
	users = models.Users{
		models.User{Id: 1},
		models.User{Id: 2},
	}
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, users)
}
