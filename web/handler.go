package web

import (
	"html/template"
	"net/http"
	"net/url"

	"github.com/POGO-Supertext-Capstone/Supertext/internal/database"
)

// V1_index handles default HTTP file requests, as well as directing to
// the home page when requested
func V1_index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.ServeFile(w, r, "web/static/index.html")
		return
	}
}

// V1_list_shows displays all shows in the database
func V1_list_shows(w http.ResponseWriter, r *http.Request) {
	db := database.UseDB(r)

	shows, err := db.GetAllShows()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render(w, "list-shows.tmpl.html", struct {
		Shows []string
	}{
		Shows: shows,
	})

}

type slide struct {
	Name    string
	ID      int
	Index   int
	HTML    template.HTML
	Content string
}

// V1_edit_show_get displays all shows in the database
func V1_edit_show_get(w http.ResponseWriter, r *http.Request) {
	showName := r.PathValue("name")
	if showName == "" {
		http.Error(w, "Invalid Show Name", http.StatusBadRequest)
		return
	}
	showName, err := url.PathUnescape(showName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Invalid Slide Index", http.StatusBadRequest)
		return
	}

}

// V1_view_show displays all shows in the database
func V1_view_show(w http.ResponseWriter, r *http.Request) {
	showName := r.PathValue("name")
	if showName == "" {
		http.Error(w, "Invalid Show Name", http.StatusBadRequest)
		return
	}

	showName, err := url.PathUnescape(showName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db := database.UseDB(r)

	data, err := db.GetShowSlides(showName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	content := make([]slide, len(data))

	for i := range data {
		content[i] = slide{
			Index:   data[i].Index,
			ID:      data[i].ID,
			HTML:    template.HTML(data[i].Content),
			Content: data[i].Content,
		}
	}

	render(w, "view-show.tmpl.html", struct {
		Name   string
		Slides []slide
	}{
		Name:   showName,
		Slides: content,
	})

}
