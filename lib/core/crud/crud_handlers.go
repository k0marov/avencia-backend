package crud

import (
	"encoding/json"
	"net/http"

	"github.com/AvenciaLab/avencia-api-contract/api/client_errors"
	"github.com/AvenciaLab/avencia-backend/lib/core/helpers/http_helpers"
	"github.com/go-chi/chi/v5"
)

type CRUDHandlers[E Entity] struct {
	store CRUDStore[E]
}

func NewCRUDEndpoint[E Entity](store CRUDStore[E]) CRUDHandlers[E] {
	return CRUDHandlers[E]{store: store}
}

func decode[E any](w http.ResponseWriter, r *http.Request) (e E, ok bool) {
	err := json.NewDecoder(r.Body).Decode(&e) 
	if err != nil {
		http_helpers.ThrowHTTPError(w, err)
		return e, false
	}
	return e, true
}

func getIdFromURL(w http.ResponseWriter, r *http.Request) (id string, ok bool) {
	id = chi.URLParam(r, "id")
	if id == "" {
		http_helpers.ThrowClientError(w, client_errors.IdNotProvided)
		return "", false
	}
	return id, true
}

func (ep CRUDHandlers[E]) Create(w http.ResponseWriter, r *http.Request) {
	e, ok := decode[E](w, r) 
	if !ok {
		return 
	}
	err := ep.store.Create(e)
	if err != nil {
		http_helpers.ThrowHTTPError(w, err)
	}
}

func (h CRUDHandlers[E]) Read(w http.ResponseWriter, r *http.Request) {
	id, ok := getIdFromURL(w, r)
	if !ok {
		return
	}
	e, err := h.store.Read(id)
	if err != nil {
		http_helpers.ThrowHTTPError(w, err)
	}
	err = json.NewEncoder(w).Encode(e)
	if err != nil {
		http_helpers.ThrowHTTPError(w, err)
	}
}

func (h CRUDHandlers[E]) Update(w http.ResponseWriter, r *http.Request) {
	e, ok := decode[E](w, r)
	if !ok {
		return
	}
	err := h.store.Update(e)
	if err != nil {
		http_helpers.ThrowHTTPError(w, err)
	}
}


func (h CRUDHandlers[E]) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := getIdFromURL(w, r)
	if !ok {
		return 
	}
	err := h.store.Delete(id)
	if err != nil {
		http_helpers.ThrowHTTPError(w, err)
	}
}

