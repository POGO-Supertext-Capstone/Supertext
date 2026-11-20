package web

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strconv"

	"github.com/POGO-Supertext-Capstone/Supertext/internal/database"
	"github.com/POGO-Supertext-Capstone/Supertext/internal/upload"
)

// V1_index handles default HTTP file requests, as well as directing to
// the home page when requested
func V1_index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.ServeFile(w, r, "web/static/index.html")
		return
	}
	http.NotFound(w, r)
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

func V1_slide_get(w http.ResponseWriter, r *http.Request) {
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

	idRaw := r.PathValue("id")
	if idRaw == "" {
		http.Error(w, "Invalid Slide Index", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idRaw)
	if err != nil {
		http.Error(w, "Invalid Slide Index", http.StatusBadRequest)
		return
	}

	db := database.UseDB(r)

	slides, err := db.GetShowSlides(showName)
	if err != nil {
		http.Error(w, fmt.Sprintf("unable to read slide content: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	var targetSlide upload.SlideContent
	for _, slide := range slides {
		if slide.ID == id {
			targetSlide = slide
			break
		}
	}
	if targetSlide == (upload.SlideContent{}) {
		http.Error(w, fmt.Sprintf("unable to read slide content: no ID found for %d", id), http.StatusInternalServerError)
		return
	}

	render(w, "slide-content.tmpl.html", slide{
		Name:    showName,
		ID:      targetSlide.ID,
		Index:   targetSlide.Index,
		HTML:    template.HTML(targetSlide.Content),
		Content: targetSlide.Content,
	})
}

func V1_slide_put(w http.ResponseWriter, r *http.Request) {
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

	idRaw := r.PathValue("id")
	if idRaw == "" {
		http.Error(w, "Invalid Slide Index", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idRaw)
	if err != nil {
		http.Error(w, "Invalid Slide Index", http.StatusBadRequest)
		return
	}

	db := database.UseDB(r)

	slides, err := db.GetShowSlides(showName)
	if err != nil {
		http.Error(w, fmt.Sprintf("unable to read slide content: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	var targetSlide upload.SlideContent
	var index int
	for i, slide := range slides {
		if slide.ID == id {
			targetSlide = slide
			index = i
			break
		}
	}
	if targetSlide == (upload.SlideContent{}) {
		http.Error(w, fmt.Sprintf("unable to read slide content: no ID found for %d", id), http.StatusInternalServerError)
		return
	}

	err = r.ParseForm()
	if err != nil {
		http.Error(w, fmt.Sprintf("unable to parse form data: %s", err.Error()), http.StatusBadRequest)
		return
	}

	newSlideContent := r.Form.Get("content")
	slides[index].Content = newSlideContent
	targetSlide.Content = newSlideContent

	err = db.UpdateShowSlides(showName, slides)
	if err != nil {
		http.Error(w, fmt.Sprintf("unable to update show content: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	render(w, "slide-content.tmpl.html", slide{
		Name:    showName,
		ID:      targetSlide.ID,
		Index:   targetSlide.Index,
		HTML:    template.HTML(targetSlide.Content),
		Content: targetSlide.Content,
	})
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

	idRaw := r.PathValue("id")
	if idRaw == "" {
		http.Error(w, "Invalid Slide Index", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idRaw)
	if err != nil {
		http.Error(w, "Invalid Slide Index", http.StatusBadRequest)
		return
	}

	db := database.UseDB(r)

	slides, err := db.GetShowSlides(showName)
	if err != nil {
		http.Error(w, fmt.Sprintf("unable to read slide content: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	var targetSlide upload.SlideContent
	for _, slide := range slides {
		if slide.ID == id {
			targetSlide = slide
			break
		}
	}
	if targetSlide == (upload.SlideContent{}) {
		http.Error(w, fmt.Sprintf("unable to read slide content: no ID found for %d", id), http.StatusInternalServerError)
		return
	}

	render(w, "edit-slide.tmpl.html", slide{
		Name:    showName,
		ID:      targetSlide.ID,
		Index:   targetSlide.Index,
		HTML:    template.HTML(targetSlide.Content),
		Content: targetSlide.Content,
	})

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
