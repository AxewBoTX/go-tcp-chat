package lib

import (
	"encoding/json"
	"net"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
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
				log.Error("Failed To Broadcast Message", "Error", message_broadcast_err)
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
		styles := log.DefaultStyles()
		styles.Levels[log.InfoLevel] = lipgloss.NewStyle().
			SetString("JOIN").
			Padding(0, 1, 0, 1).
			Background(lipgloss.Color("#a6e3a1")).
			Foreground(lipgloss.Color("0"))
		logger := log.NewWithOptions(os.Stdout, log.Options{
			ReportTimestamp: true,
		})
		logger.SetStyles(styles)
		logger.Info(msg.Client.Username + " joined the server!")
	} else if msg.Method == "MSG" {
		styles := log.DefaultStyles()
		styles.Levels[log.InfoLevel] = lipgloss.NewStyle().
			SetString("MSG").
			Padding(0, 1, 0, 1).
			Background(lipgloss.Color("#74c7ec")).
			Foreground(lipgloss.Color("0"))
		logger := log.NewWithOptions(os.Stdout, log.Options{
			ReportTimestamp: true,
		})
		logger.SetStyles(styles)
		logger.Info("", "username", msg.Client.Username, "body", msg.Body)
	} else if msg.Method == "LEAVE" {
		styles := log.DefaultStyles()
		styles.Levels[log.InfoLevel] = lipgloss.NewStyle().
			SetString("LEAVE").
			Padding(0, 1, 0, 1).
			Background(lipgloss.Color("#f9e2af")).
			Foreground(lipgloss.Color("0"))
		logger := log.NewWithOptions(os.Stdout, log.Options{
			ReportTimestamp: true,
		})
		logger.SetStyles(styles)
		logger.Info(msg.Client.Username + " left the server!")
	}
}
