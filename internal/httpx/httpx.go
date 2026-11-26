package httpx

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/nerinejaponica/api/internal/constant"
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
	// SET RESPONSE HEADER AND STATUS
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	// WHEN WRITING RESPONSE RETURNS ERROR
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("httpx: gagal menulis respons", "err", err)
	}
}

func Error(w http.ResponseWriter, status int, code, message string, details any) {
	JSON(w, status, errorEnvelope{Error: ErrorBody{Code: code, Message: message, Details: details}})
}

// LOG THE ERROR DETAIL ONLY, NEVER SEND IT TO THE CLIENT
func InternalError(w http.ResponseWriter, r *http.Request, err error) {
	slog.Error("internal error", "method", r.Method, "path", r.URL.Path, "err", err)
	Error(w, http.StatusInternalServerError, constant.CodeInternalError, constant.MsgInternalError, nil)
}

func DecodeJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	// LIMIT REQUEST BODY SIZE
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
	// REJECT UNKNOWN FIELDS
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	// WHEN DECODE BODY RETURNS ERROR
	if err := dec.Decode(dst); err != nil {
		var maxErr *http.MaxBytesError
		switch {
		case errors.Is(err, io.EOF):
			return errors.New(constant.MsgBodyEmpty)
		case errors.As(err, &maxErr):
			return fmt.Errorf(constant.MsgBodyTooLargeFmt, maxErr.Limit)
		default:
			return fmt.Errorf(constant.MsgBodyInvalidJSON, err)
		}
	}
	// REJECT BODY THAT CONTAINS MORE THAN ONE JSON OBJECT
	if dec.More() {
		return errors.New(constant.MsgBodySingleObject)
	}
	return nil
}
