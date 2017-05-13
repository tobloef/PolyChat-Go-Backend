package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func openDatabase() error {
	file, _ := os.Open("mysql_config.json")
	decoder := json.NewDecoder(file)
	config := Config{}
	err := decoder.Decode(&config)
	if err != nil {
		return fmt.Errorf("Error loading database config file\n%v", err)
	}
	dsn := fmt.Sprintf("%v:%v@/%v", config.User, config.Password, config.Database)
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("Error opening connection to database\n%v", err)
	}
	db.SetMaxIdleConns(100)
	err = db.Ping()
	if err != nil {
		return fmt.Errorf("Error opening connection to database\n%v", err)
	}
	return nil
}

func insertMessage(userId int, content string) error {
	if db == nil {
		return fmt.Errorf("Error inserting message. db is nil")
	}
	statement, err := db.Prepare("INSERT INTO messages (user_id, content) VALUES (?, ?)")
	if err != nil {
		return fmt.Errorf("Error inserting message. Couldn't prepare query\n%v", err)
	}
	defer statement.Close()
	_, err = statement.Exec(userId, content)
	if err != nil {
		return fmt.Errorf("Error inserting message. Couldn't execute query\n%v", err)
	}
	return nil
}

func insertUser(nickname string) (int, error) {
	if db == nil {
		return 0, fmt.Errorf("Error inserting user. db is nil")
	}
	statement, err := db.Prepare("INSERT INTO users (nickname) VALUES (?)")
	if err != nil {
		return 0, fmt.Errorf("Error inserting user. Couldn't prepare query\n%v", err)
	}
	defer statement.Close()
	res, err := statement.Exec(nickname)
	if err != nil {
		return 0, fmt.Errorf("Error inserting user. Couldn't execute query\n%v", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("Error inserting user. Couldn't get inserted id\n%v", err)
	}
	return int(id), nil
}

func getMessages(amount int) ([]Message, error) {
	if db == nil {
		return nil, fmt.Errorf("Error getting messages. db is nil")
	}
	var query = "SELECT users.nickname, messages.content FROM messages INNER JOIN users ON messages.user_id=users.id ORDER BY messages.id DESC"
	if amount > 0 {
		query += " LIMIT ?"
	}
	statement, err := db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("Error getting messages. Couldn't prepare query\n%v", err)
	}
	defer statement.Close()

	var rows *sql.Rows
	if amount > 0 {
		rows, err = statement.Query(amount)
	} else {
		rows, err = statement.Query()
	}
	if err != nil {
		return nil, fmt.Errorf("Error getting messages. Couldn't execute query\n%v", err)
	}
	defer rows.Close()
	messages := []Message{}
	for rows.Next() {
		var message Message
		err := rows.Scan(&message.Nickname, &message.Content)
		if err == nil {
			messages = append(messages, message)
		}
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("Error getting messages. Couldn't execute query\n%v", err)
	}
	return messages, nil
}
