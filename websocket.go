package main

import (
	"github.com/gorilla/websocket"
	"net/http"
	"encoding/json"
	"reflect"
	"fmt"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}
var clients = make(map[*websocket.Conn]Client)

func connect(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, w.Header())
	if err != nil {
		return
	}
	defer closeConnection(conn)
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			break
		}
		if reflect.DeepEqual(message, []byte("ping")) {
			err := conn.WriteMessage(messageType, []byte("pong"))
			if err != nil {
				fmt.Printf("Error responding to ping\n%v\n", err)
			}
			continue
		}
		var event Event
		err = json.Unmarshal(message, &event)
		if (err != nil) {
			fmt.Printf("Error unmarshaling event\n%v\n", err)
			continue
		}
		data, ok := event.Data.(string)
		if !ok {
			fmt.Printf("Error converting event data to string\n")
			continue
		}
		switch event.Type {
			case "connect":
				connectEvent(conn, data)
			case "message":
				if client, ok := clients[conn]; ok {
					chatMessage := Message{
						client.Nickname,
						data,
					}
					messageEvent(conn, chatMessage)
				}
		}
	}
}

func connectEvent(conn *websocket.Conn, nickname string) {
	if !nicknameAvailable(nickname) {
		event := Event{
			"connectResponse",
			"nicknameTaken",
		}
		sendEvent(conn, event)
		return
	}
	id, err := insertUser(nickname)
	if (err != nil) {
		event := Event{
			"connectResponse",
			"error",
		}
		sendEvent(conn, event)
		return
	}
	clients[conn] = Client{
		nickname,
		id,
	}
	fmt.Printf("%v connected\n", nickname)
	event := Event{
		"connectResponse",
		"ready",
	}
	sendEvent(conn, event)
	event = Event{
		"onlineCount",
		len(clients),
	}
	for connKey := range clients {
		sendEvent(connKey, event)
	}
	event = Event{
		"connected",
		nickname,
	}
	for connKey := range clients {
		if (conn != connKey) {
			sendEvent(connKey, event)
		}
	}
}

func messageEvent(conn *websocket.Conn, message Message) {
	event := Event{
		"connected",
		message,
	}
	for connKey := range clients {
		if (conn != connKey) {
			sendEvent(connKey, event)
		}
	}
	insertMessage(message)
}

func nicknameAvailable(nickname string) (bool) {
	for _, client := range clients {
		if (client.Nickname == nickname) {
			return false
		}
	}
	return true
}

func closeConnection(conn *websocket.Conn) {
	defer conn.Close()
	if client, ok := clients[conn]; ok {
		delete(clients, conn)
		event := Event{
			"disconnected",
			client.Nickname,
		}
		fmt.Printf("%v disconnected\n", client.Nickname)
		for connKey := range clients {
			sendEvent(connKey, event)
		}
		event = Event{
			"onlineCount",
			len(clients),
		}
		for connKey := range clients {
			sendEvent(connKey, event)
		}
	}
}

func sendEvent(conn *websocket.Conn, event Event) {
	eventJson, err := json.Marshal(event)
	if err != nil {
		fmt.Printf("Error marshaling event\n%v\n", err)
		return;
	}
	err = conn.WriteMessage(websocket.TextMessage, eventJson)
	if (err != nil) {
		fmt.Printf("Error writing message\n%v\n", err)
		return
	}
}