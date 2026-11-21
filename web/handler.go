package web

import (
	"net/http"
)

// V1_index handles default HTTP file requests, as well as directing to
// the home page when requested
func V1_index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Redirect(w, r, portal_url, http.StatusFound)
		return

	}

	http.ServeFile(w, r, "web/static/index.html")
}
