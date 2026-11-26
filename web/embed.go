// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package web

import (
	"embed"
	"net/http"
)

//go:embed static js
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
