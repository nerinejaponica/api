package validate

import (
	"fmt"
	"regexp"

	"github.com/nerinejaponica/api/internal/constant"
)

const (
	SlugMaxLen = 200
)

var slugPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

func MsgMaxLen(n int) string {
	return fmt.Sprintf(constant.MsgMaxLenFmt, n)
}

func SlugError(slug string) string {
	// RETURN THE FIRST RULE THAT SLUG VIOLATES
	switch {
	case slug == "":
		return constant.MsgFieldRequired
	case len(slug) > SlugMaxLen:
		return MsgMaxLen(SlugMaxLen)
	case !slugPattern.MatchString(slug):
		return constant.MsgSlugFormat
	}
	// EMPTY STRING MEANS SLUG IS VALID
	return ""
}
