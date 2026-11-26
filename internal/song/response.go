package song

import (
	"time"

	"github.com/nerinejaponica/api/internal/db"
)

type Response struct {
	ID          int64      `json:"id"`
	Slug        string     `json:"slug"`
	TitleJA     string     `json:"title_ja"`
	Status      string     `json:"status"`
	PublishedAt *time.Time `json:"published_at"`
}

// WRAPPER FOR LIST RESPONSE
type listResponse struct {
	Data []Response `json:"data"`
}

// MAP DATABASE MODEL INTO API RESPONSE
func toResponse(s db.Song) Response {
	return Response{
		ID:          s.ID,
		Slug:        s.Slug,
		TitleJA:     s.TitleJa,
		Status:      s.Status,
		PublishedAt: s.PublishedAt,
	}
}
