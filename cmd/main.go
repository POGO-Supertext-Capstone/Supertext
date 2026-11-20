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
	"log"
	"net/http"
	"os"

	"github.com/POGO-Supertext-Capstone/Supertext/internal/database"
	"github.com/POGO-Supertext-Capstone/Supertext/web"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	_, devMode := os.LookupEnv("DEV")

	ds, err := database.Open()
	if err != nil {
		log.Printf("Failed to open DataStore %s", err)
		return
	}
	defer ds.Close()

	web.UseDevTemplates(devMode)

	http.Handle("GET /v1/shows", database.HandlerFunc(ds, web.V1_list_shows))
	http.Handle("GET /v1/show/{name}/{id}/edit", database.HandlerFunc(ds, web.V1_edit_show_get))
	http.Handle("GET /v1/show/{name}/{id}", database.HandlerFunc(ds, web.V1_slide_get))
	http.Handle("PUT /v1/show/{name}/{id}", database.HandlerFunc(ds, web.V1_slide_put))
	http.Handle("GET /v1/show/{name}", database.HandlerFunc(ds, web.V1_view_show))
	http.Handle("GET /v1/upload", database.HandlerFunc(ds, web.V1_upload_get))
	http.Handle("POST /v1/upload", database.HandlerFunc(ds, web.V1_upload_post))
	http.Handle("/", database.HandlerFunc(ds, web.V1_index))

	http.Handle("/web/", web.StaticHandler(devMode))

	log.Fatal(http.ListenAndServe(":8081", nil))
}
