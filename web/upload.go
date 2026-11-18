package web

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/POGO-Supertext-Capstone/Supertext/internal/database"
	"github.com/POGO-Supertext-Capstone/Supertext/internal/upload"
)

// V1_upload_get handles all GET requests coming through the
// /v1/upload endpoint. Returns the HTML page with the form to upload
// the PPTX file
func V1_upload_get(w http.ResponseWriter, r *http.Request) {
	p := "web/static/upload.html"
	http.ServeFile(w, r, p)
}

// V1_upload_post handles all POST requests coming through the
// /v1/upload endpoint. Decodes the file passed through the form
// using the upload package
func V1_upload_post(w http.ResponseWriter, r *http.Request) {
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

	ds := database.UseDB(r)

	// The goroutine has an independent copy of the file data.
	go uploadPPTX(ds, name, h.Filename, bytes.NewReader(data))

	fmt.Fprintf(w, "<p>File Uploaded</p>")
}

// uploadPPTX handles the decoding and uplloading of a PPTX file
// into the appropiate format to be saved in the database, and
// parsed by other parts of the program
func uploadPPTX(ds *database.DataStore, name string, filename string, r io.Reader) {
	slides, err := upload.DecodePPTX(r)
	if err != nil {
		log.Printf("Failed to decode PPTX %q: %v", filename, err)
		return
	}

	ds.UploadShow(slides, name, filename)
}
