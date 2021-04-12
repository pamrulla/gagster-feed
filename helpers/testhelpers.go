package helpers

import (
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func CreateNewRouter() *chi.Mux {
	router := chi.NewRouter()
	router.Use(
		render.SetContentType(render.ContentTypeJSON),
		middleware.Logger,
		middleware.RedirectSlashes,
		middleware.Recoverer,
	)
	return router
}

func RunRequest(method string, ts *httptest.Server, urlPart string, body io.Reader) (*http.Response, error) {
	req, _ := http.NewRequest(method, ts.URL+urlPart, body)
	return http.DefaultClient.Do(req)
}
