package song

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nerinejaponica/api/internal/httpx"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.list)
	r.Post("/", h.create)
	return r
}

type listResponse struct {
	Data []Response `json:"data"`
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	songs, err := h.svc.List(r.Context())
	if err != nil {
		httpx.InternalError(w, r, err)
		return
	}
	httpx.JSON(w, http.StatusOK, listResponse{Data: songs})
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var req CreateRequest
	if err := httpx.DecodeJSON(w, r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid_json", err.Error(), nil)
		return
	}

	req.Normalize()
	if errs := req.Validate(); len(errs) > 0 {
		httpx.Error(w, http.StatusUnprocessableEntity, "validation_failed", "Data lagu tidak valid.", errs)
		return
	}

	song, err := h.svc.Create(r.Context(), req)
	switch {
	case errors.Is(err, ErrSlugTaken):
		httpx.Error(w, http.StatusConflict, "slug_taken", "Slug sudah dipakai lagu lain.", map[string]string{"slug": req.Slug})
		return
	case err != nil:
		httpx.InternalError(w, r, err)
		return
	}
	httpx.JSON(w, http.StatusCreated, song)
}
