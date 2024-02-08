package lib

import (
	"encoding/json"
	"log"
	"net"

	"github.com/fatih/color"
)

const ( // Constants
	SERVER_ADDR string = "localhost:8080"
)

type ( // Objects
	Client struct {
		Conn     net.Conn `json:"-"`
		Username string   `json:"username"`
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

// broadcast message to the network
func BroadcastMessage(msg Message) {
	for _, client := range Clients {
		if client.Conn.RemoteAddr().String() != msg.Client.Conn.RemoteAddr().String() {
			msg_broadcast_encoder := json.NewEncoder(client.Conn)
			if message_broadcast_err := msg_broadcast_encoder.Encode(msg); message_broadcast_err != nil {
				color.Set(color.FgRed)
				log.Println("Message Broadcast Error:", message_broadcast_err)
				color.Unset()
			}
		}
	}
	LogMsg(msg)
}

// generating body for message
func GenerateMsgBody(msg Message) string {
	if msg.Method == "JOIN" {
		return msg.Client.Username + " joined the server!"
	} else if msg.Method == "MSG" {
		return msg.Client.Username + ": " + msg.Body
	} else if msg.Method == "LEAVE" {
		return msg.Client.Username + " left the server!"
	} else {
		return ""
	}
}

// logging message to the server
func LogMsg(msg Message) {
	if msg.Method == "JOIN" {
		color.Set(color.FgGreen)
		log.Println(msg.Client.Username + " joined the server!")
		color.Unset()
	} else if msg.Method == "MSG" {
		log.Println(msg.Client.Username + ": " + msg.Body)
	} else if msg.Method == "LEAVE" {
		color.Set(color.FgYellow)
		log.Println(msg.Client.Username + " left the server!")
		color.Unset()
	}
}
