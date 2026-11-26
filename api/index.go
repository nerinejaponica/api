// VERCEL ENTRY POINT, vercel.json ROUTES ALL /api/* TO Handler
package handler

import (
	"context"
	"log/slog"
	"net/http"
	"sync"

	"github.com/nerinejaponica/api/internal/constant"
	"github.com/nerinejaponica/api/internal/httpx"
	"github.com/nerinejaponica/api/internal/server"
)

// KEEP AS GLOBAL SO THE POOL IS REUSED WHILE THE FUNCTION INSTANCE IS WARM
var (
	once    sync.Once
	srv     *server.Server
	initErr error
)

func Handler(w http.ResponseWriter, r *http.Request) {
	// INIT SERVER ONLY ONCE PER FUNCTION INSTANCE
	once.Do(func() {
		srv, initErr = server.New(context.Background())
	})
	// WHEN INIT SERVER RETURNS ERROR
	if initErr != nil {
		slog.Error("gagal inisialisasi server", "err", initErr)
		httpx.Error(w, http.StatusInternalServerError, constant.CodeServerInitFailed, constant.MsgServerInitFailed, nil)
		return
	}
	// HAND OVER REQUEST TO SERVER
	srv.ServeHTTP(w, r)
}
