package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/POGO-Supertext-Capstone/Supertext/internal/upload"
)

// v1_index handles default HTTP file requests, as well as directing to
// the home page when requested
func v1_index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.ServeFile(w, r, "web/static/index.html")
		return
	}
	if strings.HasPrefix(r.URL.Path, "/web/") {
		http.ServeFile(w, r, "."+r.URL.Path)
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

	http.Error(w, "Failed to read uploaded file", http.StatusBadRequest)
	if err != nil {
		return
	}
	defer file.Close()

	// create temp file to hold pptx data
	tmpFile, err := os.CreateTemp("", "presentation-*"+path.Ext(h.Filename))
	if err != nil {
		http.Error(w, "Failed to create temporary file", http.StatusInternalServerError)
		return
	}
	defer tmpFile.Close()

	// copy contents to temp file
	if _, err := io.Copy(tmpFile, file); err != nil {
		http.Error(w, "Failed to save uploaded file", http.StatusInternalServerError)
		return
	}

	// collect data for presentation upload
	var uploadData upload.PresentationUpload
	uploadData.PresentationName = name      // Presentation name (provided by user)
	uploadData.Filename = h.Filename        // Original file name (provided by user)
	uploadData.TmpFilePath = tmpFile.Name() // location of temp file

	// Decode PPTX
	go uploadPPTX(uploadData)
}

// uploadPPTX handles the decoding and uplloading of a PPTX file
// into the appropiate format to be saved in the database, and
// parsed by other parts of the program
func uploadPPTX(uploadData upload.PresentationUpload) {
	slides, err := upload.DecodePPTX(uploadData)
	if err != nil {
		log.Printf("Failed to decode PPTX: %v", err)
		return
	}
	upload.DebugPrintSlides(slides)
	// upload to database (TBI)
}

func main() {
	http.HandleFunc("/", v1_index)
	http.HandleFunc("GET /v1/upload", v1_upload_get)
	http.HandleFunc("POST /v1/upload", v1_upload_post)

	log.Fatal(http.ListenAndServe(":8081", nil))
}
