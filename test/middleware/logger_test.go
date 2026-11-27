package middleware_test

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/nerinejaponica/api/internal/middleware"
)

// captureLogs REDIRECTS slog INTO A BUFFER FOR THE DURATION OF THE TEST
func captureLogs(t *testing.T) *bytes.Buffer {
	t.Helper()
	var buf bytes.Buffer
	previous := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, nil)))
	t.Cleanup(func() { slog.SetDefault(previous) })
	return &buf
}

func TestLogger_LogsMethodPathAndStatus(t *testing.T) {
	logs := captureLogs(t)
	handler := middleware.Logger(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}))

	handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/api/songs", nil))

	for _, want := range []string{"method=POST", "path=/api/songs", "status=201"} {
		if !strings.Contains(logs.String(), want) {
			t.Errorf("log = %q, want it to contain %q", logs.String(), want)
		}
	}
}

func TestLogger_DefaultsToStatusOKWhenHandlerWritesNothing(t *testing.T) {
	logs := captureLogs(t)
	handler := middleware.Logger(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))

	if !strings.Contains(logs.String(), "status=200") {
		t.Errorf("log = %q, want status=200", logs.String())
	}
}

func TestLogger_PassesResponseThroughUnchanged(t *testing.T) {
	captureLogs(t)
	handler := middleware.Logger(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	}))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	if rec.Code != http.StatusTeapot {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusTeapot)
	}
}
