package lib

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

type Server struct {
	listener net.Listener
}

func NewServer() *Server {
	listener, listener_err := net.Listen("tcp", SERVER_ADDR)
	if listener_err != nil {
		log.Fatal("Failed To Start Server", "Error", listener_err)
	}
	return &Server{listener: listener}
}

func (server *Server) Start() {
	log.Info("Server running", "URL", "http://"+SERVER_ADDR)

	defer func() {
		server.listener.Close()
		styles := log.DefaultStyles()
		styles.Levels[log.InfoLevel] = lipgloss.NewStyle().
			Padding(0, 1, 0, 1).
			Background(lipgloss.Color("#f9e2af")).
			Foreground(lipgloss.Color("0"))
		logger := log.New(os.Stdout)
		logger.SetStyles(styles)
		logger.Info("Server closing on", "https://"+SERVER_ADDR)
	}()

	go server.ManageConnections()

	for {
		conn, conn_err := server.listener.Accept()
		if conn_err != nil {
			log.Fatal("Failed To Accept Connection", "Error", conn_err)
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
		log.Error("Failed To Decode JOIN Message", "Error", username_msg_decode_err)
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
				log.Error("Failed To Decode Client Message", "Error", client_msg_decode_err)
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
