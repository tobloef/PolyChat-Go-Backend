package main

import (
	"net/url"
	"net/http"
	"fmt"
	"encoding/json"
	"strconv"
)

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

func messages(w http.ResponseWriter, r *http.Request) {
	params, err := url.ParseQuery(r.URL.RawQuery)
	if (err != nil) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid query parameters"))
		return
	}
	var amount int
	if (len(params["amount"]) > 0) {
		amount, err = strconv.Atoi(params["amount"][0])
		if (err != nil) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Invalid amount: %v", params["amount"][0])))
			return
		}
	}
	messages, err := getMessages(amount)
	if (err != nil) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Server error: Couldn't get messages"))
		return
	}
	messagesJson, err := json.Marshal(messages)
	if (err != nil) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Server error: Couldn't convert messages to JSON"))
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write([]byte(messagesJson))
}