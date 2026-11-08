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
	"context"
	"net/http"
)

const handler_key = "supertext_publisher_context_key"

// Handler enhances the request object to allow for use of the UseSlidePublisher()
// function to access the SlidePublisher in later requests
// Designed to allow taking in a http.Handler
func Handler(sp *SlidePublisher, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		withHandler := context.WithValue(r.Context(), handler_key, sp)
		next.ServeHTTP(w, r.WithContext(withHandler))
	})
}

// HandlerFunc enhances the request object to allow for use of the UseSlidePublisher()
// function to access the SlidePublisher in later requests
// Designed to allow taking in a http.HandlerFunc
func HandlerFunc(sp *SlidePublisher, next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		withHandler := context.WithValue(r.Context(), handler_key, sp)
		next.ServeHTTP(w, r.WithContext(withHandler))
	})
}

// UseSlidePublisher pulls the SlidePublisher out of a request if the request object
// has been enhanced with the SlidePublisher
func UseSlidePublisher(r *http.Request) *SlidePublisher {
	store := r.Context().Value(handler_key)
	return store.(*SlidePublisher)
}
