package main

import (
	"bufio"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/nerinejaponica/api/internal/server"
)

func main() {
	// LOAD ".env" FILE INTO ENVIRONMENT VARIABLES
	if err := loadDotEnv(".env"); err != nil {
		slog.Error("gagal membaca .env", "err", err)
		os.Exit(1)
	}

	// INIT CONTEXT THAT IS CANCELLED ON INTERRUPT OR SIGTERM
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// INIT SERVER
	srv, err := server.New(ctx)
	// WHEN INIT SERVER RETURNS ERROR
	if err != nil {
		slog.Error("gagal inisialisasi server", "err", err)
		os.Exit(1)
	}
	defer srv.Close()

	// GET PORT FROM ENVIRONMENT WITH 8080 AS DEFAULT
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	// INIT HTTP SERVER
	httpServer := &http.Server{
		Addr:              ":" + port,
		Handler:           srv,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// RUN HTTP SERVER IN BACKGROUND
	go func() {
		slog.Info("API lokal jalan", "url", "http://localhost:"+port)
		// WHEN RUNNING SERVER RETURNS ERROR OTHER THAN GRACEFUL CLOSE
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server berhenti", "err", err)
			stop()
		}
	}()

	// WAIT UNTIL INTERRUPT OR SIGTERM IS RECEIVED
	<-ctx.Done()
	slog.Info("mematikan server...")
	// SHUTDOWN HTTP SERVER GRACEFULLY WITH 5 SECONDS LIMIT
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		slog.Error("shutdown gagal", "err", err)
	}
}

func loadDotEnv(path string) error {
	// OPEN ".env" FILE
	f, err := os.Open(path)
	// TREAT MISSING FILE AS SUCCESS BECAUSE ENVIRONMENT CAN COME FROM THE SYSTEM
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	defer f.Close()

	// READ FILE LINE BY LINE
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// SKIP EMPTY LINE AND COMMENT
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// SPLIT LINE INTO KEY AND VALUE
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(strings.TrimPrefix(key, "export "))
		value = strings.Trim(strings.TrimSpace(value), `"'`)
		// SET ENVIRONMENT VARIABLE ONLY WHEN IT IS NOT SET YET
		if _, exists := os.LookupEnv(key); !exists {
			os.Setenv(key, value)
		}
	}
	return scanner.Err()
}
