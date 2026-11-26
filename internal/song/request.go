package song

import (
	"strings"

	"github.com/nerinejaponica/api/internal/constant"
	"github.com/nerinejaponica/api/internal/validate"
)

type CreateRequest struct {
	Slug    string `json:"slug"`
	TitleJA string `json:"title_ja"`
	Status  string `json:"status"`
}

func (req *CreateRequest) Normalize() {
	// TRIM WHITESPACE FROM INPUT
	req.Slug = strings.TrimSpace(req.Slug)
	req.TitleJA = strings.TrimSpace(req.TitleJA)
	req.Status = strings.TrimSpace(req.Status)
	// USE DRAFT AS DEFAULT STATUS
	if req.Status == "" {
		req.Status = StatusDraft
	}
}

func (req CreateRequest) Validate() map[string]string {
	// COLLECT ERROR PER FIELD
	errs := map[string]string{}
	// VALIDATE SLUG
	if msg := validate.SlugError(req.Slug); msg != "" {
		errs["slug"] = msg
	}
	// VALIDATE TITLE
	if req.TitleJA == "" {
		errs["title_ja"] = constant.MsgFieldRequired
	}
	// VALIDATE STATUS
	if req.Status != StatusDraft && req.Status != StatusPublished {
		errs["status"] = constant.MsgInvalidStatus
	}
	return errs
}
