// Copyright (C) 2026 Jamie Murphy
//
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

package events

import (
	"fmt"
	"net/http"
	"time"
	"bytes"
	"encoding/json"

	"github.com/POGO-Supertext-Capstone/Supertext/internal/database"
	"github.com/POGO-Supertext-Capstone/Supertext/internal/upload"
)

type Show struct {
	ShowName string `json:"show"`
	Slides []upload.SlideContent `json:"slides"`
}

func V1_Slide_Content(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	ds := database.UseDB(r)
	shows, err := ds.GetAllShows()
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}

	// FIXME: Get the show name somewhere else
	showSlides, err := ds.GetShowSlides(shows[0])

	show := Show { ShowName: shows[0], Slides: showSlides }

	// Serialize data
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	err = encoder.Encode(show)
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}
	data := buf.String()

	w.Header().Set("HX-Trigger", "slidesLoaded")
	fmt.Fprintf(w, "%s", data)
}

func V1_SSE_Viewer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	ds := database.UseDB(r)
	shows, err := ds.GetAllShows()
	if err != nil {
		return
	}

	clientDisconnect := r.Context().Done()

	rc := http.NewResponseController(w)
	t := time.NewTicker(time.Second)
	defer t.Stop()

	_, err = fmt.Fprintf(w, "event: Event\ndata: %s\n\n", shows)

	idx := 0
	for {
		select {
		case <-clientDisconnect:
			return
		case <-t.C:
			_, err := fmt.Fprintf(w, "event: Event\ndata: %d\n\n", idx)
			idx++
			if err != nil {
				return
			}
			err = rc.Flush()
			if err != nil {
				return
			}
		}
	}
}
