package lib

import (
	"net"

	"github.com/fatih/color"
)

const (
	SERVER_ADDR string = "localhost:8080"
)

type Client struct {
}

func (server *Server) HandleConnection(conn net.Conn) {
	color.Green(conn.RemoteAddr().String() + " joined the server!\n")
	defer func() {
		color.Yellow(conn.RemoteAddr().String() + " left the server!\n")
		conn.Close()
	}()
}

func (server *Server) ManageConnections() {
}
