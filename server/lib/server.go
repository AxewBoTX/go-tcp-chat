package lib

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"

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
			log.Println("Connection Accept Error:", conn_err)
			color.Unset()
		}

		go server.HandleConnection(conn)
	}
}

// Handle individual connections
func (server *Server) HandleConnection(conn net.Conn) {
	client_addr := conn.RemoteAddr().String()

	// Get Username
	if _, conn_write_err := conn.Write([]byte("Enter your username:")); conn_write_err != nil {
		log.Println("Conn Write Error:", conn_write_err)
	}
	username, username_read_err := bufio.NewReader(conn).ReadString('\n')
	if username_read_err != nil {
		fmt.Fprintln(conn, "Username Read Error:", username_read_err)
		conn.Close()
	}
	username = strings.TrimRight(strings.TrimSpace(username), "\n")
	if username == "" || len(username) == 0 {
		fmt.Fprintln(conn, "Error: Username is required")
	}

	// Create Client
	client := *NewClient(conn, username)
	Clients[client_addr] = client
	JOIN_Messages <- *NewMessage(client, "JOIN", "")

	// Read Loop
	go func() {
		client_scanner := bufio.NewScanner(conn)
		for client_scanner.Scan() {
			Messages <- *NewMessage(client, "MSG", client_scanner.Text())
		}
		defer func() {
			LEAVE_Messages <- *NewMessage(client, "LEAVE", "")
			delete(Clients, client_addr)
			conn.Close()
		}()
	}()
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
