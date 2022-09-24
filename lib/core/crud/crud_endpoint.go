package crud

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewCrudEndpoint[E Entity](h CRUDHandlers[E]) http.HandlerFunc {
	r := chi.NewRouter()
	r.Post("/", h.Create)
	r.Get("/{id}", h.Read)
	r.Patch("/{id}", h.Update)
	r.Delete("/{id}", h.Delete)
	return r.ServeHTTP
}
