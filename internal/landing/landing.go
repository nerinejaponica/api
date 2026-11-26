package landing

import (
	_ "embed"
	"net/http"
)

//go:embed index.html
var page []byte

func Handler(w http.ResponseWriter, r *http.Request) {
	// SEND EMBEDDED LANDING PAGE
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(page)
}
