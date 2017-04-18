package main

import (
	"net/http"
)

func main() {
	println("Starting server...")
	http.HandleFunc("/api/ping", ping)
	http.HandleFunc("/api/messages", messages)
	http.ListenAndServe(":3001", nil)
}
