package health

import (
	"context"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nerinejaponica/api/internal/constant"
	"github.com/nerinejaponica/api/internal/httpx"
)

type Handler struct {
	pool      *pgxpool.Pool
	startedAt time.Time
}

func NewHandler(pool *pgxpool.Pool) *Handler {
	return &Handler{pool: pool, startedAt: time.Now()}
}

type response struct {
	Status   string `json:"status"`
	Database string `json:"database"`
	Uptime   int64  `json:"uptime"`
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	// LIMIT DATABASE PING TO 3 SECONDS
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	// COUNT UPTIME SINCE HANDLER IS CREATED
	uptime := int64(time.Since(h.startedAt).Seconds())

	// WHEN PING DATABASE RETURNS ERROR
	if err := h.pool.Ping(ctx); err != nil {
		httpx.JSON(w, http.StatusServiceUnavailable, response{Status: constant.HealthDegraded, Database: constant.DatabaseDown, Uptime: uptime})
		return
	}
	// SUCCESS PING DATABASE
	httpx.JSON(w, http.StatusOK, response{Status: constant.HealthOK, Database: constant.DatabaseUp, Uptime: uptime})
}
