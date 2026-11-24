package health

import (
	"context"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nerinejaponica/api/internal/httpx"
)

type Handler struct {
	pool *pgxpool.Pool
}

func NewHandler(pool *pgxpool.Pool) *Handler {
	return &Handler{pool: pool}
}

type response struct {
	Status   string `json:"status"`
	Database string `json:"database"`
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	if err := h.pool.Ping(ctx); err != nil {
		httpx.JSON(w, http.StatusServiceUnavailable, response{Status: "degraded", Database: "down"})
		return
	}
	httpx.JSON(w, http.StatusOK, response{Status: "ok", Database: "up"})
}
