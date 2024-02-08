package lib

import (
	"bufio"
	"encoding/json"
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
	username_json, username_json_read_err := bufio.NewReader(conn).ReadString('\n')
	if username_json_read_err != nil {
		fmt.Fprintln(conn, "Username JSON Read Error:", username_json_read_err)
		conn.Close()
	}
	var username_msg Message
	if username_msg_decode_err := json.Unmarshal([]byte(username_json), &username_msg); username_msg_decode_err != nil {
		color.Set(color.FgRed)
		log.Println("JOIN Message Decode Error:", username_msg_decode_err)
		color.Unset()
	}
	username := strings.TrimSpace(username_msg.Client.Username)

	// Create Client
	client := *NewClient(conn, username)
	Clients[client_addr] = client
	JOIN_Messages <- *NewMessage(client, "JOIN", "")

	// Read Loop
	go func() {
		client_scanner := bufio.NewScanner(conn)
		for client_scanner.Scan() {
			var client_msg Message
			if client_msg_decode_err := json.Unmarshal([]byte(client_scanner.Text()), &client_msg); client_msg_decode_err != nil {
				color.Set(color.FgRed)
				log.Println("Client Message Decode Error:", client_msg_decode_err)
				color.Unset()
			}
			client_msg.Client.Conn = client.Conn
			Messages <- client_msg
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
			BroadcastMessage(data)
		case data := <-Messages: // Message Broadcaster
			BroadcastMessage(data)
		case data := <-LEAVE_Messages: // Client Leave Message
			BroadcastMessage(data)
		}
	}
}
