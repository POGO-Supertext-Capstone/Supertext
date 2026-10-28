package main

import (
	"fmt"
	"strings"
	"os"
	"io"
	"path"
	"net/http"
	"log"

	"github.com/POGO-Supertext-Capstone/Supertext/internal/upload"
)

func v1_index(w http.ResponseWriter, r *http.Request) {
  if r.URL.Path == "/" {
    http.ServeFile(w, r, "web/static/index.html")
    return
  }
	if (strings.HasPrefix(r.URL.Path, "/web/")) {
		http.ServeFile(w, r, "."+r.URL.Path)
		return
	}
}

func v1_upload(w http.ResponseWriter, r *http.Request) {
	fmt.Print("/v1/upload : ")
	switch r.Method {
	case http.MethodGet:
		fmt.Println("GET")
		
		p := "web/static/upload.html"
		http.ServeFile(w, r, p)
	case http.MethodPost:
		fmt.Println("POST")

		// grab form data
		file, h, err := r.FormFile("file")
		name := r.FormValue("name")

    if err != nil {
      http.Error(w, "Failed to read uploaded file", http.StatusBadRequest)
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
		uploadData.PresentationName = name			// Presentation name (provided by user)
		uploadData.Filename = h.Filename				// Original file name (provided by user)
		uploadData.TmpFilePath = tmpFile.Name()	// location of temp file 

		// Decode PPTX
		go uploadPPTX(uploadData);

    fmt.Println("Uploaded:", h.Filename)
    fmt.Println("Saved to:", tmpFile.Name())

	default:
		// Return 405
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func uploadPPTX(uploadData upload.PresentationUpload) {
	slides, err := upload.DecodePPTX(uploadData);
	if (err != nil) {
		log.Printf("Failed to decode PPTX: %v", err)
		return
	}
	// upload to dtabase
}

func main() {
	http.HandleFunc("/", v1_index)
	http.HandleFunc("/v1/upload", v1_upload)

	log.Fatal(http.ListenAndServe(":8081", nil))
}
