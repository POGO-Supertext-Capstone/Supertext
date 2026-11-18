package web

import (
	"embed"
	"net/http"
)

//go:embed static js templates
var embedFS embed.FS

// StaticHandler handlers file based web content. If in dev mode, the handler
// will pull from the filesystem on each request, in production it will instead
// pull from an embedded filesystem allowing for easier builds since the
// supporting files will be included in the final binary
func StaticHandler(dev bool) http.Handler {
	var fileServer http.Handler

	if dev {
		fileServer = http.FileServer(http.Dir("web/"))
	} else {
		fileServer = http.FileServer(http.FS(embedFS))
	}

	return http.StripPrefix("/web/", fileServer)
}
