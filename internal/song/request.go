package song

import (
	"strings"

	"github.com/nerinejaponica/api/internal/validate"
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
	if msg := validate.SlugError(req.Slug); msg != "" {
		errs["slug"] = msg
	}
	if req.TitleJA == "" {
		errs["title_ja"] = validate.MsgRequired
	}
	if req.Status != StatusDraft && req.Status != StatusPublished {
		errs["status"] = msgInvalidStatus
	}
	return errs
}
