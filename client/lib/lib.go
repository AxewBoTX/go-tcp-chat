package lib

import (
	"net"
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
