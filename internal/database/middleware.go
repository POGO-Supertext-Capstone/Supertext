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

package database

import (
	"context"
	"net/http"
)

const handler_key = "supertext_DB_context_key"

// Handler enhances the request object to allow for use of the UseDB() function
// to access the DataStore in later requests
// Designed to allow taking in a http.Handler
func Handler(db *DataStore, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		withDB := context.WithValue(r.Context(), handler_key, db)
		next.ServeHTTP(w, r.WithContext(withDB))
	})
}

// HandlerFunc enhances the request object to allow for use of the UseDB()
// function to access the DataStore in later requests
// Designed to allow taking in a http.HandlerFunc
func HandlerFunc(db *DataStore, next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		withDB := context.WithValue(r.Context(), handler_key, db)
		next.ServeHTTP(w, r.WithContext(withDB))
	})
}

// UseDB pulls out the DataStore out of a request if the request object has
// been enhanced with the DataStore
func UseDB(r *http.Request) *DataStore {
	store := r.Context().Value(handler_key)
	return store.(*DataStore)
}
