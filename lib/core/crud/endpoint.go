package crud

import (
	"github.com/go-chi/chi/v5"
)

type EnabledMethods struct {
	Create, Read, Update, Delete bool 
}

func NewCrudEndpoint[E Entity](h Handlers[E], m EnabledMethods) func(chi.Router) {
	return func(r chi.Router) {
		if m.Create {
			r.Post("/", h.Create)
		}
		if m.Read {
			r.Get("/{id}", h.Read)
			r.Get("/", h.Read)
		}
		if m.Update {
			r.Patch("/{id}", h.Update)
			r.Patch("/", h.Update)
		}
		if m.Delete {
			r.Delete("/{id}", h.Delete)
		}
	}
}
