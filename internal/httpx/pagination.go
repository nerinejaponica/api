package httpx

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/nerinejaponica/api/internal/constant"
)

const (
	DefaultLimit = 20
	MaxLimit     = 50
)

// Cursor MARKS THE LAST ITEM OF A PAGE. LISTS ARE SORTED BY published_at DESC, id DESC,
// SO BOTH VALUES ARE NEEDED: id BREAKS TIES BETWEEN ITEMS WITH THE SAME published_at.
// PublishedAt IS NIL FOR A DRAFT, WHICH ONLY SHOWS UP IN ADMIN LISTS.
type Cursor struct {
	PublishedAt *time.Time `json:"p"`
	ID          int64      `json:"i"`
}

// Page IS THE PAGINATION INPUT OF A LIST ENDPOINT. Cursor IS NIL FOR THE FIRST PAGE.
type Page struct {
	Limit  int
	Cursor *Cursor
}

// ListResponse IS THE SHAPE OF EVERY LIST ENDPOINT. NextCursor IS NIL ON THE LAST PAGE.
type ListResponse[T any] struct {
	Data       []T     `json:"data"`
	NextCursor *string `json:"next_cursor"`
}

func (c Cursor) Encode() string {
	// KEEP THE TIME IN UTC SO THE SAME INSTANT ALWAYS GIVES THE SAME CURSOR
	if c.PublishedAt != nil {
		utc := c.PublishedAt.UTC()
		c.PublishedAt = &utc
	}
	raw, _ := json.Marshal(c)
	// URL SAFE BASE64 WITHOUT PADDING, SO IT CAN SIT IN A QUERY STRING AS IT IS
	return base64.RawURLEncoding.EncodeToString(raw)
}

func DecodeCursor(encoded string) (Cursor, error) {
	var cursor Cursor
	invalid := errors.New(constant.MsgCursorInvalid)

	// DECODE BASE64
	raw, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return Cursor{}, invalid
	}
	// DECODE JSON AND REJECT UNKNOWN FIELDS
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&cursor); err != nil {
		return Cursor{}, invalid
	}
	// A REAL ROW ALWAYS HAS A POSITIVE ID
	if cursor.ID < 1 {
		return Cursor{}, invalid
	}
	return cursor, nil
}

// ParsePage READS ?limit= AND ?cursor=. IT RETURNS ERRORS PER FIELD, EMPTY WHEN THE INPUT IS VALID.
func ParsePage(r *http.Request) (Page, map[string]string) {
	query := r.URL.Query()
	errs := map[string]string{}

	// READ LIMIT
	limit, msg := parseLimit(query.Get("limit"))
	if msg != "" {
		errs["limit"] = msg
	}
	// READ CURSOR
	cursor, msg := parseCursor(query.Get("cursor"))
	if msg != "" {
		errs["cursor"] = msg
	}
	return Page{Limit: limit, Cursor: cursor}, errs
}

func parseLimit(raw string) (int, string) {
	// MISSING LIMIT MEANS DEFAULT
	if raw == "" {
		return DefaultLimit, ""
	}
	// LIMIT MUST BE A WHOLE NUMBER BETWEEN 1 AND MaxLimit
	limit, err := strconv.Atoi(raw)
	if err != nil || limit < 1 || limit > MaxLimit {
		return DefaultLimit, fmt.Sprintf(constant.MsgLimitInvalidFmt, MaxLimit)
	}
	return limit, ""
}

func parseCursor(raw string) (*Cursor, string) {
	// MISSING CURSOR MEANS FIRST PAGE
	if raw == "" {
		return nil, ""
	}
	// DECODE CURSOR
	cursor, err := DecodeCursor(raw)
	if err != nil {
		return nil, err.Error()
	}
	return &cursor, ""
}

// NewListResponse EXPECTS THE QUERY TO FETCH LIMIT+1 ROWS. THE EXTRA ROW ONLY TELLS THAT
// ANOTHER PAGE EXISTS: IT IS CUT OFF AND THE CURSOR POINTS AT THE LAST ROW THAT STAYS.
func NewListResponse[T any](rows []T, limit int, cursorOf func(T) Cursor) ListResponse[T] {
	// NEVER RETURN NULL LIST
	if rows == nil {
		rows = []T{}
	}
	// LAST PAGE: NO EXTRA ROW, NO NEXT CURSOR
	if len(rows) <= limit {
		return ListResponse[T]{Data: rows}
	}
	// MORE PAGES: CUT THE EXTRA ROW AND POINT THE CURSOR AT THE LAST KEPT ROW
	rows = rows[:limit]
	next := cursorOf(rows[limit-1]).Encode()
	return ListResponse[T]{Data: rows, NextCursor: &next}
}
