package httpx_test

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/nerinejaponica/api/internal/constant"
	"github.com/nerinejaponica/api/internal/httpx"
)

type item struct {
	ID          int64
	PublishedAt *time.Time
}

func cursorOfItem(i item) httpx.Cursor {
	return httpx.Cursor{PublishedAt: i.PublishedAt, ID: i.ID}
}

// requestWith BUILDS A GET REQUEST WITH THE GIVEN QUERY PARAMETERS
func requestWith(params map[string]string) *http.Request {
	values := url.Values{}
	for key, value := range params {
		values.Set(key, value)
	}
	return httptest.NewRequest(http.MethodGet, "/api/songs?"+values.Encode(), nil)
}

// rawCursor ENCODES ANY TEXT AS A CURSOR, TO BUILD MALFORMED ONES
func rawCursor(text string) string {
	return base64.RawURLEncoding.EncodeToString([]byte(text))
}

func TestParsePage_NoParamsUsesDefaultLimitAndNoCursor(t *testing.T) {
	page, errs := httpx.ParsePage(requestWith(nil))

	if len(errs) != 0 {
		t.Fatalf("errs = %v, want none", errs)
	}
	if page.Limit != httpx.DefaultLimit || page.Cursor != nil {
		t.Errorf("page = %+v, want default limit and nil cursor", page)
	}
}

func TestParsePage_ValidLimitIsUsed(t *testing.T) {
	page, errs := httpx.ParsePage(requestWith(map[string]string{"limit": "5"}))

	if len(errs) != 0 || page.Limit != 5 {
		t.Errorf("limit = %d, errs = %v, want 5 and no errors", page.Limit, errs)
	}
}

func TestParsePage_LimitAtMaxIsAccepted(t *testing.T) {
	page, errs := httpx.ParsePage(requestWith(map[string]string{"limit": "50"}))

	if len(errs) != 0 || page.Limit != httpx.MaxLimit {
		t.Errorf("limit = %d, errs = %v, want max limit and no errors", page.Limit, errs)
	}
}

func TestParsePage_LimitOverMaxIsRejected(t *testing.T) {
	_, errs := httpx.ParsePage(requestWith(map[string]string{"limit": "51"}))

	if errs["limit"] == "" {
		t.Errorf("errs = %v, want an error for limit", errs)
	}
}

func TestParsePage_LimitZeroIsRejected(t *testing.T) {
	_, errs := httpx.ParsePage(requestWith(map[string]string{"limit": "0"}))

	if errs["limit"] == "" {
		t.Errorf("errs = %v, want an error for limit", errs)
	}
}

func TestParsePage_NegativeLimitIsRejected(t *testing.T) {
	_, errs := httpx.ParsePage(requestWith(map[string]string{"limit": "-3"}))

	if errs["limit"] == "" {
		t.Errorf("errs = %v, want an error for limit", errs)
	}
}

func TestParsePage_LimitNotANumberIsRejected(t *testing.T) {
	_, errs := httpx.ParsePage(requestWith(map[string]string{"limit": "banyak"}))

	if errs["limit"] == "" {
		t.Errorf("errs = %v, want an error for limit", errs)
	}
}

func TestParsePage_LimitErrorMessageMentionsMax(t *testing.T) {
	_, errs := httpx.ParsePage(requestWith(map[string]string{"limit": "999"}))

	if !strings.Contains(errs["limit"], "50") {
		t.Errorf("message = %q, want it to mention the max limit 50", errs["limit"])
	}
}

func TestParsePage_ValidCursorIsDecoded(t *testing.T) {
	publishedAt := time.Date(2026, 11, 1, 10, 0, 0, 0, time.UTC)
	encoded := httpx.Cursor{PublishedAt: &publishedAt, ID: 42}.Encode()

	page, errs := httpx.ParsePage(requestWith(map[string]string{"cursor": encoded}))

	if len(errs) != 0 || page.Cursor == nil || page.Cursor.ID != 42 {
		t.Errorf("cursor = %+v, errs = %v, want id 42 and no errors", page.Cursor, errs)
	}
}

func TestParsePage_InvalidCursorIsRejected(t *testing.T) {
	page, errs := httpx.ParsePage(requestWith(map[string]string{"cursor": "bukan-cursor"}))

	if errs["cursor"] != constant.MsgCursorInvalid || page.Cursor != nil {
		t.Errorf("cursor = %+v, errs = %v, want invalid cursor error", page.Cursor, errs)
	}
}

func TestParsePage_ReportsLimitAndCursorErrorsTogether(t *testing.T) {
	_, errs := httpx.ParsePage(requestWith(map[string]string{"limit": "0", "cursor": "xx"}))

	if errs["limit"] == "" || errs["cursor"] == "" {
		t.Errorf("errs = %v, want errors for both limit and cursor", errs)
	}
}

func TestCursor_RoundTripKeepsPublishedAtAndID(t *testing.T) {
	publishedAt := time.Date(2026, 11, 1, 10, 30, 15, 123456000, time.UTC)

	got, err := httpx.DecodeCursor(httpx.Cursor{PublishedAt: &publishedAt, ID: 7}.Encode())

	if err != nil {
		t.Fatalf("err = %v, want nil", err)
	}
	if got.ID != 7 || got.PublishedAt == nil || !got.PublishedAt.Equal(publishedAt) {
		t.Errorf("got %+v, want id 7 and published_at %v", got, publishedAt)
	}
}

func TestCursor_RoundTripOfDraftKeepsNilPublishedAt(t *testing.T) {
	got, err := httpx.DecodeCursor(httpx.Cursor{ID: 9}.Encode())

	if err != nil || got.PublishedAt != nil || got.ID != 9 {
		t.Errorf("got %+v, err %v, want nil published_at and id 9", got, err)
	}
}

func TestCursor_SameInstantInOtherZoneGivesSameCursor(t *testing.T) {
	utc := time.Date(2026, 11, 1, 3, 0, 0, 0, time.UTC)
	jakarta := utc.In(time.FixedZone("WIB", 7*3600))

	a := httpx.Cursor{PublishedAt: &utc, ID: 1}.Encode()
	b := httpx.Cursor{PublishedAt: &jakarta, ID: 1}.Encode()

	if a != b {
		t.Errorf("cursors differ for the same instant: %q vs %q", a, b)
	}
}

func TestCursor_EncodedValueIsSafeInQueryString(t *testing.T) {
	publishedAt := time.Date(2026, 11, 1, 10, 0, 0, 0, time.UTC)

	encoded := httpx.Cursor{PublishedAt: &publishedAt, ID: 1234567}.Encode()

	if strings.ContainsAny(encoded, "+/=") {
		t.Errorf("cursor %q contains characters that need escaping", encoded)
	}
}

func TestDecodeCursor_NotBase64IsRejected(t *testing.T) {
	_, err := httpx.DecodeCursor("***")

	if err == nil {
		t.Error("err = nil, want an error")
	}
}

func TestDecodeCursor_NotJSONIsRejected(t *testing.T) {
	_, err := httpx.DecodeCursor(rawCursor("bukan json"))

	if err == nil {
		t.Error("err = nil, want an error")
	}
}

func TestDecodeCursor_UnknownFieldIsRejected(t *testing.T) {
	_, err := httpx.DecodeCursor(rawCursor(`{"p":null,"i":1,"x":2}`))

	if err == nil {
		t.Error("err = nil, want an error")
	}
}

func TestDecodeCursor_ZeroIDIsRejected(t *testing.T) {
	_, err := httpx.DecodeCursor(rawCursor(`{"p":null,"i":0}`))

	if err == nil {
		t.Error("err = nil, want an error")
	}
}

func TestDecodeCursor_EmptyInputIsRejected(t *testing.T) {
	_, err := httpx.DecodeCursor("")

	if err == nil {
		t.Error("err = nil, want an error")
	}
}

func TestNewListResponse_MoreRowsThanLimitCutsExtraAndSetsNextCursor(t *testing.T) {
	first := time.Date(2026, 11, 3, 0, 0, 0, 0, time.UTC)
	second := time.Date(2026, 11, 2, 0, 0, 0, 0, time.UTC)
	third := time.Date(2026, 11, 1, 0, 0, 0, 0, time.UTC)
	rows := []item{{ID: 3, PublishedAt: &first}, {ID: 2, PublishedAt: &second}, {ID: 1, PublishedAt: &third}}

	got := httpx.NewListResponse(rows, 2, cursorOfItem)

	if len(got.Data) != 2 || got.Data[1].ID != 2 {
		t.Errorf("data = %+v, want the first two rows", got.Data)
	}
	if got.NextCursor == nil {
		t.Fatal("next_cursor = nil, want a cursor")
	}
}

func TestNewListResponse_NextCursorPointsAtLastKeptRow(t *testing.T) {
	first := time.Date(2026, 11, 3, 0, 0, 0, 0, time.UTC)
	second := time.Date(2026, 11, 2, 0, 0, 0, 0, time.UTC)
	rows := []item{{ID: 3, PublishedAt: &first}, {ID: 2, PublishedAt: &second}, {ID: 1}}

	got := httpx.NewListResponse(rows, 2, cursorOfItem)
	cursor, err := httpx.DecodeCursor(*got.NextCursor)

	if err != nil || cursor.ID != 2 || !cursor.PublishedAt.Equal(second) {
		t.Errorf("cursor = %+v, err %v, want the second row (id 2)", cursor, err)
	}
}

func TestNewListResponse_RowsEqualToLimitHaveNoNextCursor(t *testing.T) {
	rows := []item{{ID: 2}, {ID: 1}}

	got := httpx.NewListResponse(rows, 2, cursorOfItem)

	if len(got.Data) != 2 || got.NextCursor != nil {
		t.Errorf("got %+v, want two rows and no next cursor", got)
	}
}

func TestNewListResponse_FewerRowsThanLimitHaveNoNextCursor(t *testing.T) {
	got := httpx.NewListResponse([]item{{ID: 1}}, 20, cursorOfItem)

	if len(got.Data) != 1 || got.NextCursor != nil {
		t.Errorf("got %+v, want one row and no next cursor", got)
	}
}

func TestNewListResponse_NoRowsGivesEmptyDataAndNullCursorInJSON(t *testing.T) {
	got := httpx.NewListResponse[item](nil, 20, cursorOfItem)

	raw, err := json.Marshal(got)

	if err != nil || string(raw) != `{"data":[],"next_cursor":null}` {
		t.Errorf("json = %s, err %v, want data [] and next_cursor null", raw, err)
	}
}
