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
	"github.com/nerinejaponica/api/internal/constant"
	"github.com/nerinejaponica/api/internal/db"
	"github.com/nerinejaponica/api/internal/health"
	"github.com/nerinejaponica/api/internal/httpx"
	"github.com/nerinejaponica/api/internal/landing"
	"github.com/nerinejaponica/api/internal/middleware"
	"github.com/nerinejaponica/api/internal/song"
)

const maxConns = 4

type Server struct {
	router http.Handler
	pool   *pgxpool.Pool
}

func New(ctx context.Context) (*Server, error) {
	// LOAD CONFIGURATION
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	// INIT DATABASE POOL
	pool, err := newPool(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	// INIT HANDLERS AND THEIR DEPENDENCIES
	queries := db.New(pool)
	healthHandler := health.NewHandler(pool)
	songHandler := song.NewHandler(song.NewService(queries))

	// INIT ROUTER AND GLOBAL MIDDLEWARE
	r := chi.NewRouter()
	r.Use(middleware.Recover, middleware.Logger)

	// SET FALLBACK RESPONSE FOR UNKNOWN ENDPOINT AND METHOD
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		httpx.Error(w, http.StatusNotFound, constant.CodeNotFound, constant.MsgRouteNotFound, nil)
	})
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		httpx.Error(w, http.StatusMethodNotAllowed, constant.CodeMethodNotAllowed, constant.MsgMethodNotAllowed, nil)
	})

	// SET ROUTES
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
	// PARSE DATABASE URL INTO POOL CONFIGURATION
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		// DO NOT INCLUDE THE PARSE ERROR BECAUSE IT CAN CONTAIN THE DATABASE PASSWORD
		return nil, errors.New(constant.MsgDatabaseURLInvalid)
	}
	// LIMIT POOL SIZE
	cfg.MaxConns = maxConns
	// USE EXEC MODE BECAUSE NEON POOLED (PGBOUNCER TRANSACTION MODE) DOES NOT SUPPORT NAMED PREPARED STATEMENTS
	cfg.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeExec

	// CREATE POOL
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf(constant.MsgPoolCreateFailed, err)
	}
	return pool, nil
}
