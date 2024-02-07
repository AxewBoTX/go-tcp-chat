package lib

import (
	"net"
)

const ( // Constants
	SERVER_ADDR string = "localhost:8080"
)

type ( // Objects
	Client struct {
		Conn     net.Conn
		Username string `json:"username"`
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

// Handle individual connections
func (server *Server) HandleConnection(conn net.Conn) {
	client_addr := conn.RemoteAddr().String()
	Clients[client_addr] = Client{Conn: conn}

	JOIN_Messages <- Message{Client_ADDR: client_addr, Body: client_addr + " joined the server"}

	go func() {
		client_scanner := bufio.NewScanner(conn)
		for client_scanner.Scan() {
			Messages <- Message{Client_ADDR: client_addr, Body: client_scanner.Text()}
		}
		defer func() {
			LEAVE_Messages <- Message{Client_ADDR: client_addr, Body: client_addr + " left the server."}
			delete(Clients, client_addr)
			conn.Close()
		}()
	}()
}

// Manage Connections
func (server *Server) ManageConnections() {
	for {
		select {
		case data := <-JOIN_Messages: // Client Join Message
			color.Set(color.FgGreen)
			for _, client := range Clients {
				if client.Conn.RemoteAddr().String() != data.Client_ADDR {
					if _, broadcast_err := client.Conn.Write([]byte(data.Client_ADDR + ": " + data.Body)); broadcast_err != nil {
						log.Println(data.Body)
					}
				}
			}
			log.Println(data.Body)
			color.Unset()
		case data := <-Messages: // Broadcaster
			for _, client := range Clients {
				if client.Conn.RemoteAddr().String() != data.Client_ADDR {
					if _, broadcast_err := client.Conn.Write([]byte(data.Client_ADDR + ": " + data.Body)); broadcast_err != nil {
						color.Set(color.FgRed)
						log.Println(data.Body)
						color.Unset()
					}
				}
			}
			fmt.Printf("-> %s\n", data.Client_ADDR+": "+data.Body)
		case data := <-LEAVE_Messages: // Client Leave Message
			color.Set(color.FgYellow)
			for _, client := range Clients {
				if client.Conn.RemoteAddr().String() != data.Client_ADDR {
					if _, broadcast_err := client.Conn.Write([]byte(data.Client_ADDR + ": " + data.Body)); broadcast_err != nil {
						log.Println(data.Body)
					}
				}
			}
			log.Println(data.Body)
			color.Unset()
		}
	}
}
