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
func (server *Server) ManageConnections() {
	for {
		select {
		case data := <-JOIN_Messages: // Client Join Message
			for _, client := range Clients {
				if client.Conn.RemoteAddr().String() != data.Client.Conn.RemoteAddr().String() {
					if _, join_message_err := client.Conn.Write([]byte(data.Client.Username + " joined the server!")); join_message_err != nil {
						color.Set(color.FgRed)
						log.Println("JOIN Message Error:", join_message_err)
						color.Unset()
					}
				}
			}
			color.Set(color.FgGreen)
			log.Println(data.Client.Username + " joined the server!")
			color.Unset()
		case data := <-Messages: // Message Broadcaster
			for _, client := range Clients {
				if client.Conn.RemoteAddr().String() != data.Client.Conn.RemoteAddr().String() {
					if _, message_broadcast_err := client.Conn.Write([]byte(data.Client.Username + ": " + data.Body)); message_broadcast_err != nil {
						color.Set(color.FgRed)
						log.Println("Message Broadcast Error:", message_broadcast_err)
						color.Unset()
					}
				}
			}
			log.Println(data.Client.Username + ": " + data.Body)
		case data := <-LEAVE_Messages: // Client Leave Message
			for _, client := range Clients {
				if client.Conn.RemoteAddr().String() != data.Client.Conn.RemoteAddr().String() {
					if _, leave_message_err := client.Conn.Write([]byte(data.Client.Username + " left the server!")); leave_message_err != nil {
						color.Set(color.FgRed)
						log.Println(leave_message_err)
						color.Unset()
					}
				}
			}
			color.Set(color.FgYellow)
			log.Println(data.Client.Username + " left the server!")
			color.Unset()
		}
	}
}
