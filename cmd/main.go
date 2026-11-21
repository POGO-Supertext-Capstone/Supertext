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
	if err != nil && !os.IsNotExist(err) {
		log.Fatalf("Error loading .env file: %s", err)
	}
	_, devMode := os.LookupEnv("DEV")

	ds, err := database.Open()
	if err != nil {
		log.Fatalf("Failed to open DataStore: %s", err)
	}
	defer ds.Close()

	http.Handle("/", database.HandlerFunc(ds, web.V1_index))
	http.Handle("GET /v1/upload", database.HandlerFunc(ds, web.V1_upload_get))
	http.Handle("POST /v1/upload", database.HandlerFunc(ds, web.V1_upload_post))

	// connectivity checks from phones and laptops that join the hotspot
	http.HandleFunc("/generate_204", web.Captive_portal_redirect)
	http.HandleFunc("/hotspot-detect.html", web.Captive_portal_redirect)
	http.HandleFunc("/connecttest.txt", web.Captive_portal_redirect)

	http.Handle("/web/", web.StaticHandler(devMode))

	log.Fatal(http.ListenAndServe(":8081", nil))
}
