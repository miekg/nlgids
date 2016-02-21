package nlgids

import (
	"io"
	"net/http"
)

// Download will present a binary download to the user.
func Download(rd io.Reader, filename string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
	io.Copy(w, rd)
}
