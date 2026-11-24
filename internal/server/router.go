package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nerinejaponica/api/internal/config"
	"github.com/nerinejaponica/api/internal/db"
	"github.com/nerinejaponica/api/internal/health"
	"github.com/nerinejaponica/api/internal/httpx"
	"github.com/nerinejaponica/api/internal/middleware"
	"github.com/nerinejaponica/api/internal/song"
	"github.com/nerinejaponica/api/internal/landing"
)

const maxConns = 4

type Server struct {
	router http.Handler
	pool   *pgxpool.Pool
}

func New(ctx context.Context) (*Server, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	pool, err := newPool(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	queries := db.New(pool)
	healthHandler := health.NewHandler(pool)
	songHandler := song.NewHandler(song.NewService(queries))

	r := chi.NewRouter()
	r.Use(middleware.Recover, middleware.Logger)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		httpx.Error(w, http.StatusNotFound, "not_found", "Endpoint tidak ditemukan.", nil)
	})
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		httpx.Error(w, http.StatusMethodNotAllowed, "method_not_allowed", "Method tidak didukung untuk endpoint ini.", nil)
	})

	r.Get("/", landing.Handler)
	r.Route("/api", func(r chi.Router) {
		r.Get("/health", healthHandler.Get)
		r.Mount("/songs", songHandler.Routes())
	})

	return &Server{router: r, pool: pool}, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) Close() {
	s.pool.Close()
}

func newPool(ctx context.Context, url string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		// err dari ParseConfig bisa memuat password, jadi tidak disertakan.
		return nil, errors.New("server: DATABASE_URL tidak bisa dibaca")
	}
	cfg.MaxConns = maxConns
	// Neon pooled = PgBouncer mode transaksi; mode Exec tidak memakai prepared statement bernama.
	cfg.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeExec

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("server: gagal membuat pool database: %w", err)
	}
	return pool, nil
}
