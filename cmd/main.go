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

package main

import (
	"bytes"
	"io"
	"log"
	"net/http"

	"github.com/POGO-Supertext-Capstone/Supertext/internal/upload"
)

// v1_index handles default HTTP file requests, as well as directing to
// the home page when requested
func v1_index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.ServeFile(w, r, "web/static/index.html")
		return
	}
}

// v1_upload_get handles all GET requests coming through the
// /v1/upload endpoint. Returns the HTML page with the form to upload
// the PPTX file
func v1_upload_get(w http.ResponseWriter, r *http.Request) {
	p := "web/static/upload.html"
	http.ServeFile(w, r, p)
}

// v1_upload_post handles all POST requests coming through the
// /v1/upload endpoint. Decodes the file passed through the form
// using the upload package
func v1_upload_post(w http.ResponseWriter, r *http.Request) {
	// grab form data
	name := r.FormValue("name")
	file, h, err := r.FormFile("file")

	if err != nil {
		http.Error(w, "Failed to read uploaded file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read the uploaded file before starting the goroutine.
	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read uploaded file", http.StatusInternalServerError)
		return
	}

	// The goroutine has an independent copy of the file data.
	go uploadPPTX(name, h.Filename, bytes.NewReader(data))
}

// uploadPPTX handles the decoding and uplloading of a PPTX file
// into the appropiate format to be saved in the database, and
// parsed by other parts of the program
func uploadPPTX(name string, filename string, r io.Reader) {
	slides, err := upload.DecodePPTX(r)
	if err != nil {
		log.Printf("Failed to decode PPTX %q: %v", filename, err)
		return
	}

	upload.DebugPrintSlides(slides)
	// upload to database (TBI)
}

func main() {
	http.HandleFunc("/", v1_index)
	http.HandleFunc("GET /v1/upload", v1_upload_get)
	http.HandleFunc("POST /v1/upload", v1_upload_post)

	fileServer := http.FileServer(http.Dir("web/"))
	http.Handle("/web/", http.StripPrefix("/web/", fileServer))

	log.Fatal(http.ListenAndServe(":8081", nil))
}
