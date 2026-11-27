package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nerinejaponica/api/internal/middleware"
)

const allowedOrigin = "http://localhost:3000"

// serve SENDS ONE REQUEST THROUGH THE CORS MIDDLEWARE, OPTIONAL ORIGIN HEADER
func serve(method, origin string) *httptest.ResponseRecorder {
	handler := middleware.CORS([]string{allowedOrigin})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(method, "/api/songs", nil)
	if origin != "" {
		req.Header.Set("Origin", origin)
	}
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return rec
}

func TestCORS_AllowedOriginGetsHeader(t *testing.T) {
	rec := serve(http.MethodGet, allowedOrigin)

	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != allowedOrigin {
		t.Errorf("allow-origin = %q, want %q", got, allowedOrigin)
	}
	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestCORS_UnknownOriginGetsNoHeader(t *testing.T) {
	rec := serve(http.MethodGet, "http://evil.example")

	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Errorf("allow-origin = %q, want empty", got)
	}
}

func TestCORS_NoOriginGetsNoHeader(t *testing.T) {
	rec := serve(http.MethodGet, "")

	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Errorf("allow-origin = %q, want empty", got)
	}
}

func TestCORS_PreflightFromAllowedOriginIsAnsweredWithNoContent(t *testing.T) {
	rec := serve(http.MethodOptions, allowedOrigin)

	if rec.Code != http.StatusNoContent {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusNoContent)
	}
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != allowedOrigin {
		t.Errorf("allow-origin = %q, want %q", got, allowedOrigin)
	}
	if got := rec.Header().Get("Access-Control-Allow-Methods"); got == "" {
		t.Error("allow-methods is empty, want the allowed methods")
	}
}

func TestCORS_PreflightFromUnknownOriginReachesRouter(t *testing.T) {
	rec := serve(http.MethodOptions, "http://evil.example")

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Errorf("allow-origin = %q, want empty", got)
	}
}
