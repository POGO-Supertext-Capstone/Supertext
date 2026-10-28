package main

import (
	"fmt"
	"strings"
	"os"
	"io"
	"path"
	"github.com/POGO-Supertext-Capstone/Supertext/internal/upload"
	"net/http"
	"log"
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
		
		file, h, err := r.FormFile("file")
    if err != nil {
        http.Error(w, "Failed to read uploaded file", http.StatusBadRequest)
      return
  	}
    defer file.Close()

    tmpFile, err := os.CreateTemp("", "presentation-*"+path.Ext(h.Filename))
    if err != nil {
      http.Error(w, "Failed to create temporary file", http.StatusInternalServerError)
      return
    }
    defer tmpFile.Close()
  	defer os.Remove(tmpFile.Name())

    if _, err := io.Copy(tmpFile, file); err != nil {
      http.Error(w, "Failed to save uploaded file", http.StatusInternalServerError)
    	return
  	}

    fmt.Println("Uploaded:", h.Filename)
    fmt.Println("Saved to:", tmpFile.Name())

	default:
		// Return 405
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func main() {
	fmt.Println("Hello Beavers!")
	upload.Test()

	http.HandleFunc("/", v1_index)
	http.HandleFunc("/v1/upload", v1_upload)

	log.Fatal(http.ListenAndServe(":8081", nil))
}
