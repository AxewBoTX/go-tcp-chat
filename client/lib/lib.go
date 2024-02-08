package lib

import (
	"fmt"
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

func PrintMSG(msg Message) {
	if msg.Method == "JOIN" {
		color.Set(color.FgGreen)
		fmt.Println(msg.Client.Username + " joined the server!")
		color.Unset()
	} else if msg.Method == "MSG" {
		fmt.Println(msg.Client.Username + ": " + msg.Body)
	} else if msg.Method == "LEAVE" {
		color.Set(color.FgYellow)
		fmt.Println(msg.Client.Username + " left the server!")
		color.Unset()
	}
}
