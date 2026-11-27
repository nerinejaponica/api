package middleware_test

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/nerinejaponica/api/internal/constant"
	"github.com/nerinejaponica/api/internal/middleware"
)

// silenceLogs KEEPS EXPECTED ERROR LOGS OUT OF THE TEST OUTPUT
func silenceLogs(t *testing.T) {
	t.Helper()
	previous := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&strings.Builder{}, nil)))
	t.Cleanup(func() { slog.SetDefault(previous) })
}

func TestRecover_PanicBecomesInternalError(t *testing.T) {
	silenceLogs(t)
	handler := middleware.Recover(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("boom")
	}))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/songs", nil))

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
	}
	if !strings.Contains(rec.Body.String(), constant.CodeInternalError) {
		t.Errorf("body = %s, want code %s", rec.Body.String(), constant.CodeInternalError)
	}
}

func TestRecover_PanicValueIsNotSentToClient(t *testing.T) {
	silenceLogs(t)
	handler := middleware.Recover(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("secret-detail")
	}))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	if strings.Contains(rec.Body.String(), "secret-detail") {
		t.Errorf("body leaks the panic value: %s", rec.Body.String())
	}
}

func TestRecover_NoPanicPassesThrough(t *testing.T) {
	handler := middleware.Recover(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	}))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	if rec.Code != http.StatusTeapot {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusTeapot)
	}
}

func TestRecover_AbortHandlerPanicIsRethrown(t *testing.T) {
	handler := middleware.Recover(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic(http.ErrAbortHandler)
	}))

	defer func() {
		if got := recover(); got != http.ErrAbortHandler {
			t.Errorf("recovered %v, want http.ErrAbortHandler", got)
		}
	}()
	handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
}
