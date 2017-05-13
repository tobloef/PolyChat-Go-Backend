package main

type Message struct {
	Nickname string `json:"nickname"`
	Content  string `json:"content"`
}

type Event struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type Client struct {
	Nickname string
	Id       int
}

type Config struct {
	Database string
	User     string
	Password string
}
