package httpx

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

const maxBodyBytes = 1024 * 1024

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

type errorEnvelope struct {
	Error ErrorBody `json:"error"`
}

func JSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("httpx: gagal menulis respons", "err", err)
	}
}

func Error(w http.ResponseWriter, status int, code, message string, details any) {
	JSON(w, status, errorEnvelope{Error: ErrorBody{Code: code, Message: message, Details: details}})
}

// Detail error hanya masuk log, tidak dikirim ke klien.
func InternalError(w http.ResponseWriter, r *http.Request, err error) {
	slog.Error("internal error", "method", r.Method, "path", r.URL.Path, "err", err)
	Error(w, http.StatusInternalServerError, "internal_error", "Terjadi kesalahan di server.", nil)
}

func DecodeJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		var maxErr *http.MaxBytesError
		switch {
		case errors.Is(err, io.EOF):
			return errors.New("body request kosong")
		case errors.As(err, &maxErr):
			return fmt.Errorf("body request melebihi %d byte", maxErr.Limit)
		default:
			return fmt.Errorf("JSON tidak valid: %w", err)
		}
	}
	if dec.More() {
		return errors.New("body request hanya boleh berisi satu objek JSON")
	}
	return nil
}
