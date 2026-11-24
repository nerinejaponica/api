package song

import (
	"regexp"
	"strings"
	"time"

	"github.com/nerinejaponica/api/internal/db"
)

var slugPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

const (
	StatusDraft     = "draft"
	StatusPublished = "published"
)

type CreateRequest struct {
	Slug    string `json:"slug"`
	TitleJA string `json:"title_ja"`
	Status  string `json:"status"`
}

func (req *CreateRequest) Normalize() {
	req.Slug = strings.TrimSpace(req.Slug)
	req.TitleJA = strings.TrimSpace(req.TitleJA)
	req.Status = strings.TrimSpace(req.Status)
	if req.Status == "" {
		req.Status = StatusDraft
	}
}

func (req CreateRequest) Validate() map[string]string {
	errs := map[string]string{}
	switch {
	case req.Slug == "":
		errs["slug"] = "wajib diisi"
	case len(req.Slug) > 200:
		errs["slug"] = "maksimal 200 karakter"
	case !slugPattern.MatchString(req.Slug):
		errs["slug"] = "hanya huruf kecil, angka, dan tanda hubung (contoh: izayoi-lionheart)"
	}
	if req.TitleJA == "" {
		errs["title_ja"] = "wajib diisi"
	}
	if req.Status != StatusDraft && req.Status != StatusPublished {
		errs["status"] = "harus draft atau published"
	}
	return errs
}

type Response struct {
	ID          int64      `json:"id"`
	Slug        string     `json:"slug"`
	TitleJA     string     `json:"title_ja"`
	Status      string     `json:"status"`
	PublishedAt *time.Time `json:"published_at"`
}

func toResponse(s db.Song) Response {
	return Response{
		ID:          s.ID,
		Slug:        s.Slug,
		TitleJA:     s.TitleJa,
		Status:      s.Status,
		PublishedAt: s.PublishedAt,
	}
}
