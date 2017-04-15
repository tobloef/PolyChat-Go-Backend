package main

import (
	"net/http"
	"net/url"
	"io/ioutil"
	"fmt"
	"strconv"
)

func main() {
	println("Starting server...")
	http.HandleFunc("/ping", ping)
	http.HandleFunc("/messages", messages)
	http.ListenAndServe(":8080", nil)
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

func messages(w http.ResponseWriter, r *http.Request) {
	params, err := url.ParseQuery(r.URL.RawQuery)
	if (err != nil) {
		w.Write([]byte("Error :("))
		return
	}
	amount, err := strconv.Atoi(params["amount"][0])
	if (err != nil) {
		w.Write([]byte("Error :("))
		return
	}
	json, err := fetchMessages(amount)
	if (err != nil) {
		w.Write([]byte("Error :("))
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write([]byte(json))
}

func fetchMessages(amount int) (string, error) {
	apiUrl := "http://tobloef.com/polychat/node/api/messages"
	if (amount > 0) {
		apiUrl += fmt.Sprintf("?amount=%d", amount)
	}
	res, err := http.Get(apiUrl)
	if (err != nil) {
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if (err != nil) {
		return "", err
	}
	return string(body), nil
}
