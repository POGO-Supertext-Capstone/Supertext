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
