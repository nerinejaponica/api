package httpx_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/nerinejaponica/api/internal/constant"
	"github.com/nerinejaponica/api/internal/httpx"
)

type payload struct {
	Name string `json:"name"`
}

// decode RUNS DecodeJSON AGAINST A RAW REQUEST BODY
func decode(body string) (payload, error) {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	var dst payload
	err := httpx.DecodeJSON(httptest.NewRecorder(), req, &dst)
	return dst, err
}

func TestJSON_WritesStatusHeaderAndBody(t *testing.T) {
	rec := httptest.NewRecorder()

	httpx.JSON(rec, http.StatusCreated, payload{Name: "izayoi"})

	if rec.Code != http.StatusCreated {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusCreated)
	}
	if got := rec.Header().Get("Content-Type"); got != "application/json; charset=utf-8" {
		t.Errorf("content-type = %q, want json utf-8", got)
	}
	var body payload
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil || body.Name != "izayoi" {
		t.Errorf("body = %q, want name izayoi (err %v)", rec.Body.String(), err)
	}
}

func TestJSON_UnencodableValueDoesNotPanic(t *testing.T) {
	rec := httptest.NewRecorder()

	httpx.JSON(rec, http.StatusOK, make(chan int))

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestError_WritesEnvelopeWithDetails(t *testing.T) {
	rec := httptest.NewRecorder()

	httpx.Error(rec, http.StatusUnprocessableEntity, constant.CodeValidationFailed, "invalid", map[string]string{"slug": "wajib diisi"})

	want := `{"error":{"code":"validation_failed","message":"invalid","details":{"slug":"wajib diisi"}}}`
	if got := strings.TrimSpace(rec.Body.String()); got != want {
		t.Errorf("body = %s, want %s", got, want)
	}
	if rec.Code != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusUnprocessableEntity)
	}
}

func TestError_OmitsDetailsWhenNil(t *testing.T) {
	rec := httptest.NewRecorder()

	httpx.Error(rec, http.StatusNotFound, constant.CodeNotFound, "missing", nil)

	if strings.Contains(rec.Body.String(), "details") {
		t.Errorf("body = %s, want no details field", rec.Body.String())
	}
}

func TestInternalError_HidesCauseFromClient(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/songs", nil)

	httpx.InternalError(rec, req, errors.New("password=secret db exploded"))

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
	}
	if strings.Contains(rec.Body.String(), "secret") {
		t.Errorf("body leaks the cause: %s", rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), constant.CodeInternalError) {
		t.Errorf("body = %s, want code %s", rec.Body.String(), constant.CodeInternalError)
	}
}

func TestDecodeJSON_ValidBody(t *testing.T) {
	got, err := decode(`{"name":"izayoi"}`)

	if err != nil || got.Name != "izayoi" {
		t.Errorf("got %+v, err %v, want name izayoi and no error", got, err)
	}
}

func TestDecodeJSON_EmptyBody(t *testing.T) {
	_, err := decode("")

	if err == nil || err.Error() != constant.MsgBodyEmpty {
		t.Errorf("err = %v, want %q", err, constant.MsgBodyEmpty)
	}
}

func TestDecodeJSON_MalformedJSON(t *testing.T) {
	_, err := decode(`{"name":`)

	if err == nil {
		t.Error("err = nil, want error for malformed json")
	}
}

func TestDecodeJSON_UnknownFieldRejected(t *testing.T) {
	_, err := decode(`{"name":"a","extra":1}`)

	if err == nil {
		t.Error("err = nil, want error for unknown field")
	}
}

func TestDecodeJSON_MoreThanOneObjectRejected(t *testing.T) {
	_, err := decode(`{"name":"a"}{"name":"b"}`)

	if err == nil || err.Error() != constant.MsgBodySingleObject {
		t.Errorf("err = %v, want %q", err, constant.MsgBodySingleObject)
	}
}

func TestDecodeJSON_BodyOverOneMegabyteRejected(t *testing.T) {
	body := `{"name":"` + strings.Repeat("a", 1024*1024) + `"}`

	_, err := decode(body)

	if err == nil || err.Error() == constant.MsgBodyEmpty {
		t.Errorf("err = %v, want a too-large error", err)
	}
}
