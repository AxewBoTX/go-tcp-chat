package lib

import (
	"net"
)

const ( // Constants
	SERVER_ADDR string = "localhost:8080"
)

type ( // Objects
	Client struct {
		Conn     net.Conn
		Username string `json:"username"`
	}
	Message struct {
		Client Client `json:"client"`
		Method string `json:"method"`
		Body   string `json:"body"`
	}
)

var ( // Variables
	Clients        = make(map[string]Client)
	JOIN_Messages  = make(chan Message)
	Messages       = make(chan Message)
	LEAVE_Messages = make(chan Message)
)

func NewClient(conn net.Conn, username string) *Client {
	return &Client{Conn: conn, Username: username}
}

func NewMessage(client Client, method string, body string) *Message {
	return &Message{Client: client, Method: method, Body: body}
}
