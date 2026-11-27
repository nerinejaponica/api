package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/nerinejaponica/api/internal/constant"
	"github.com/nerinejaponica/api/internal/httpx"
)

func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			// GET PANIC VALUE
			rec := recover()
			if rec == nil {
				return
			}
			// LET THE SERVER HANDLE ABORTED REQUEST
			if rec == http.ErrAbortHandler {
				panic(rec)
			}
			// LOG PANIC WITH STACK TRACE AND SEND GENERIC ERROR TO CLIENT
			slog.Error("panic", "method", r.Method, "path", r.URL.Path, "panic", rec, "stack", string(debug.Stack()))
			httpx.Error(w, http.StatusInternalServerError, constant.CodeInternalError, constant.MsgInternalError, nil)
		}()
		next.ServeHTTP(w, r)
	})
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (s *statusRecorder) WriteHeader(code int) {
	s.status = code
	s.ResponseWriter.WriteHeader(code)
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		// WRAP WRITER TO CAPTURE RESPONSE STATUS
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)
		// LOG REQUEST AFTER IT IS HANDLED
		slog.Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rec.status,
			"duration_ms", time.Since(start).Milliseconds(),
		)
	})
}

// CORS ALLOWS BROWSER REQUESTS ONLY FROM THE LISTED ORIGINS
func CORS(allowedOrigins []string) func(http.Handler) http.Handler {
	// BUILD LOOKUP SET ONCE, NOT PER REQUEST
	allowed := make(map[string]struct{}, len(allowedOrigins))
	for _, origin := range allowedOrigins {
		allowed[origin] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			_, ok := allowed[origin]

			// SKIP REQUEST WITHOUT ORIGIN OR FROM UNKNOWN ORIGIN, BROWSER WILL BLOCK IT
			if origin == "" || !ok {
				next.ServeHTTP(w, r)
				return
			}

			// ALLOW THIS ORIGIN, VARY SO CACHES DO NOT MIX ORIGINS
			w.Header().Add("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Origin", origin)

			// ANSWER PREFLIGHT HERE, DO NOT PASS IT TO THE ROUTER
			if r.Method == http.MethodOptions {
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
				w.Header().Set("Access-Control-Max-Age", "600")
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
