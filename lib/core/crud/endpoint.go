package crud

import (
	"github.com/go-chi/chi/v5"
)

func NewCrudEndpoint[E Entity](h Handlers[E]) func(chi.Router) {
	return func(r chi.Router) {
		r.Post("/", h.Create)
		r.Get("/{id}", h.Read)
		r.Patch("/{id}", h.Update)
		r.Delete("/{id}", h.Delete)
	}
}
