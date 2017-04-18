package main

import (
	"fmt"
)

func open() (error) {
	fmt.Printf("Fake opening the database connection.")
	return nil
}

func close() (error) {
	fmt.Printf("Fake closing the database connection.")
	return nil
}

//							 not []string VVV      VVV not string
func executeQuery(query string, values []string) (string, error) {
	fmt.Printf("Fake executing query:\n%v", query)
	return "", nil
}

func insertMessage(message Message) (error) {
	fmt.Printf("Fake inserting message:\n%v: %v", message.Nickname, message.Content)
	return nil
}

func getMessages(amount int) ([]Message, error) {
	fmt.Printf("Fake getting the newest %v messages.", amount)
	messages := []Message{
		Message{
			"Lukas",
			"Du er mega sej :D:D:D:D:D:D XDXDXDXDXD",
		},
		Message{
			"Lukas",
			"Hej!!!!!!",
		},
		Message{
			"Tobias",
			"Hej",
		},
	}
	if (amount != 0) {
		if len(messages) < amount {
        	amount = len(messages)
		}
		messages = messages[:amount]
	}
	return messages, nil
}