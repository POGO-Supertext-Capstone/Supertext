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

package present

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/POGO-Supertext-Capstone/Supertext/internal/database"
	"github.com/POGO-Supertext-Capstone/Supertext/internal/upload"
)

// Provided to /view endpoint as containing the show details.
type Show struct {
	ShowName string                `json:"show"`
	Slides   []upload.SlideContent `json:"slides"`
}

// V1_Slides handles all requests for the /v1/slides
// endpoint. Returns a structured response containing a show/slide body.
func V1_Slides(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ds := database.UseDB(r)
	shows, err := ds.GetAllShows()
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}

	// FIXME [sprint 3]: The show ID should be provided as a query/path parameter
	slides, err := ds.GetShowSlides(shows[0])
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}

	show := Show{ShowName: shows[0], Slides: slides}

	// Serialize data
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)

	// SECURITY: We are required to provide HTML to the client for swapping. It should be expected
	// that the shows are escaped before being placed within the database.
	encoder.SetEscapeHTML(false)
	err = encoder.Encode(show)
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}
	data := buf.String()

	w.Header().Set("HX-Trigger", "slides-loaded")
	fmt.Fprintf(w, "%s", data)
}

// V1_Stream handles all requests for the /v1/stream endpoint.
// Streams events to clients for real-time slide updates.
func V1_Stream(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming is not supported on this endpoint", 500)
		return
	}

	sp := UseSlidePublisher(r)
	client := sp.Subscribe()

	defer sp.Unsubscribe(client)

	// Infinite-loop, waits for a message on the client channel or for the Request to have been
	// closed
	for {
		select {
		case <-r.Context().Done():
			return
		// When a message is received from the client, write it to the response and flush the buffer
		case msg := <-client:
			fmt.Fprintf(w, "%s\n\n", msg)
			flusher.Flush()
		}
	}
}

// V1_Present_Post handles all POST requests coming through the
// /v1/present/* endpoints. Updates the current state of the show.
func V1_Present_Post(w http.ResponseWriter, r *http.Request) {
	sp := UseSlidePublisher(r)

	if strings.Contains(r.URL.Path, "/forward") {
		sp.show = sp.show + 1
	} else if strings.Contains(r.URL.Path, "/backward") {
		sp.show = sp.show - 1
	}

	sp.Notify(fmt.Sprintf("event: slideupdate\ndata: %d", sp.show))

	fmt.Fprintf(w, "%d", sp.show)
}
