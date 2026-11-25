package validate

import (
	"fmt"
	"regexp"
)

const (
	MsgRequired   = "wajib diisi"
	MsgSlugFormat = "hanya huruf kecil, angka, dan tanda hubung (contoh: izayoi-lionheart)"

	SlugMaxLen = 200
)

var slugPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

func MsgMaxLen(n int) string {
	return fmt.Sprintf("maksimal %d karakter", n)
}

func SlugError(slug string) string {
	switch {
	case slug == "":
		return MsgRequired
	case len(slug) > SlugMaxLen:
		return MsgMaxLen(SlugMaxLen)
	case !slugPattern.MatchString(slug):
		return MsgSlugFormat
	}
	return ""
}
