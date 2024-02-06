package lib

import (
	"log"
	"net"

	"github.com/fatih/color"
)

type Server struct {
	listener net.Listener
}

func NewServer() *Server {
	listener, listener_err := net.Listen("tcp", SERVER_ADDR)
	if listener_err != nil {
		color.Set(color.FgRed)
		log.Println("Server Listener Error:", listener_err)
		color.Unset()
	}
	return &Server{listener: listener}
}

func (server *Server) Start() {
	color.Set(color.FgHiCyan)
	log.Println("Server running on", "http://"+SERVER_ADDR)
	color.Unset()

	defer func() {
		server.listener.Close()
		color.Set(color.FgYellow)
		log.Println("Server closing on", "https://"+SERVER_ADDR)
		color.Unset()
	}()

	go server.ManageConnections()

	for {
		conn, conn_err := server.listener.Accept()
		if conn_err != nil {
			color.Set(color.FgRed)
			log.Println("Conneciton Accetp Error:", conn_err)
			color.Unset()
		}

		go server.HandleConnection(conn)
	}
}
