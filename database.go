package main

import (
	"fmt"
)

func openDatabase() (error) {
	fmt.Printf("Fake opening the database connection.\n")
	return nil
}

func closeDatabase() (error) {
	fmt.Printf("Fake closing the database connection.\n")
	return nil
}

//							 not []string VVV      VVV not string
func executeQuery(query string, values []string) (string, error) {
	fmt.Printf("Fake executing query:\n%v\n", query)
	return "", nil
}

func insertMessage(message Message) (error) {
	fmt.Printf("%v: %v\n", message.Nickname, message.Content)
	return nil
}

func insertUser(nickname string) (int, error) {
	fmt.Printf("Fake inserting user with nickname %v\n", nickname)
	return 0, nil
}

func getMessages(amount int) ([]Message, error) {
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
	fmt.Printf("Fake getting the newest %v messages.\n", len(messages))
	return messages, nil
}