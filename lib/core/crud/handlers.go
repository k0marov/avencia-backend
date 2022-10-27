package crud

import (
	"encoding/json"
	"net/http"

	"github.com/AvenciaLab/avencia-backend/lib/core/helpers/http_helpers"
	"github.com/AvenciaLab/avencia-backend/lib/features/auth/domain/service"
	"github.com/go-chi/chi/v5"
)

type Handlers[E Entity] struct {
	service *Service[E]
}

func NewCRUDHandlers[E Entity](service *Service[E]) Handlers[E] {
	return Handlers[E]{service: service}
}

func decode[E any](w http.ResponseWriter, r *http.Request) (e E, ok bool) {
	err := json.NewDecoder(r.Body).Decode(&e) 
	if err != nil {
		http_helpers.ThrowHTTPError(w, err)
		return e, false
	}
	return e, true
}
func getRequestData(r *http.Request) RequestData {
	user, _ := service.UserFromCtx(r.Context())
	return RequestData{
		IdFromURL: chi.URLParam(r, "id"),
		CallerId:  user.Id,
	}
}

func (ep Handlers[E]) Create(w http.ResponseWriter, r *http.Request) {
	e, ok := decode[E](w, r) 
	if !ok {
		return 
	}
	err := ep.service.Create(getRequestData(r), e)
	if err != nil {
		http_helpers.ThrowHTTPError(w, err)
	}
}

func (h Handlers[E]) Read(w http.ResponseWriter, r *http.Request) {
	e, err := h.service.Read(getRequestData(r))
	if err != nil {
		http_helpers.ThrowHTTPError(w, err)
		return 
	}
	err = json.NewEncoder(w).Encode(e)
	if err != nil {
		http_helpers.ThrowHTTPError(w, err)
		return 
	}
}

func (h Handlers[E]) Update(w http.ResponseWriter, r *http.Request) {
	e, ok := decode[E](w, r)
	if !ok {
		return
	}
	err := h.service.Update(getRequestData(r), e)
	if err != nil {
		http_helpers.ThrowHTTPError(w, err)
		return 
	}
}


func (h Handlers[E]) Delete(w http.ResponseWriter, r *http.Request) {
	err := h.service.Delete(getRequestData(r))
	if err != nil {
		http_helpers.ThrowHTTPError(w, err)
	}
}

