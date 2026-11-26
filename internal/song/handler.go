package song

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nerinejaponica/api/internal/constant"
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

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	// GET ALL SONGS
	songs, err := h.svc.List(r.Context())
	if err != nil {
		httpx.InternalError(w, r, err)
		return
	}
	// SUCCESS GET SONGS
	httpx.JSON(w, http.StatusOK, listResponse{Data: songs})
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	// DECODE REQUEST BODY
	var req CreateRequest
	if err := httpx.DecodeJSON(w, r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, constant.CodeInvalidJSON, err.Error(), nil)
		return
	}

	// NORMALIZE AND VALIDATE REQUEST
	req.Normalize()
	if errs := req.Validate(); len(errs) > 0 {
		httpx.Error(w, http.StatusUnprocessableEntity, constant.CodeValidationFailed, constant.MsgSongInvalid, errs)
		return
	}

	// CREATE SONG
	song, err := h.svc.Create(r.Context(), req)
	switch {
	case errors.Is(err, ErrSlugTaken):
		httpx.Error(w, http.StatusConflict, constant.CodeSlugTaken, constant.MsgSlugTakenClient, map[string]string{"slug": req.Slug})
		return
	case err != nil:
		httpx.InternalError(w, r, err)
		return
	}
	// SUCCESS CREATE SONG
	httpx.JSON(w, http.StatusCreated, song)
}
