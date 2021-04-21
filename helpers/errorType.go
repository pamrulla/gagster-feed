package helpers

import (
	"net/http"

	"github.com/go-chi/render"
)

type Error struct {
	Message string `json:"message"`
}

func NewError(w http.ResponseWriter, r *http.Request, msg string, code int) {
	render.Status(r, code)
	render.JSON(w, r, Error{Message: msg})
}
