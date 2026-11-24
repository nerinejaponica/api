// Entry point Vercel; vercel.json mengarahkan semua /api/* ke Handler.
package handler

import (
	"context"
	"log/slog"
	"net/http"
	"sync"

	"github.com/nerinejaponica/api/internal/httpx"
	"github.com/nerinejaponica/api/internal/server"
)

// Global supaya pool dipakai ulang selama instance fungsi masih hangat.
var (
	once    sync.Once
	srv     *server.Server
	initErr error
)

func Handler(w http.ResponseWriter, r *http.Request) {
	once.Do(func() {
		srv, initErr = server.New(context.Background())
	})
	if initErr != nil {
		slog.Error("gagal inisialisasi server", "err", initErr)
		httpx.Error(w, http.StatusInternalServerError, "server_init_failed", "Server gagal dijalankan.", nil)
		return
	}
	srv.ServeHTTP(w, r)
}
