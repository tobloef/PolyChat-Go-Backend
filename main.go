package main

import (
	"net/http"
)

func main() {
	println("Starting server...")
	http.HandleFunc("/api/ping", defaultHeaders(ping))
	http.HandleFunc("/api/messages", defaultHeaders(messages))
	http.HandleFunc("/", connect)
	http.ListenAndServe(":3001", nil)
}

func defaultHeaders(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
	        w.Header().Set("Access-Control-Allow-Origin", origin)
	    }
	    handler(w, r)
	}
}