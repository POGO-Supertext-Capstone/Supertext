package web

import (
	"net/http"
)

// home page on the hotpost and address that the hardware repo gives the Rasperry Pi 5
const portal_url = "http://10.42.0.1:8081/"

// answers the connectivity checks that phones and laptops make when they join a network
// answering with redirect tells them the network has a sign-in page
func Captive_portal_redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, portal_url, http.StatusFound)
}
