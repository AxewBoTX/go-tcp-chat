package lib

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/fatih/color"
)

func NewClient() *Client {
	conn, conn_err := net.Dial("tcp", SERVER_ADDR)
	if conn_err != nil {
		color.Set(color.FgRed)
		log.Fatal("TCP Connection Error:", conn_err)
		color.Unset()
	}
	return &Client{Conn: conn}
}

func (client *Client) ReadLoop() {
	// Client JOIN Message
	join_msg_encoder := json.NewEncoder(client.Conn)
	join_msg := Message{Client: *client, Method: "JOIN"}
	if join_msg_encoder_err := join_msg_encoder.Encode(join_msg); join_msg_encoder_err != nil {
		color.Set(color.FgRed)
		log.Fatal("JOIN Message Encoder Error:", join_msg_encoder_err)
		color.Unset()
	}
	go func() { // Running in a goroutine
		for {
			read_buff := make([]byte, 2048)
			if _, read_err := client.Conn.Read(read_buff); read_err != nil {
				color.Set(color.FgRed)
				log.Fatal("Server Recieve Error:", read_err)
				color.Unset()
				return
			}
			fmt.Printf("-> %s\n", read_buff)
		}
	}()
}

func (client *Client) WriteLoop() {
	user_input_scanner := bufio.NewScanner(os.Stdin)
	encoder := json.NewEncoder(client.Conn)
	for user_input_scanner.Scan() {
		msg := Message{Client: *client, Method: "MSG", Body: user_input_scanner.Text()}
		if user_msg_encode_err := encoder.Encode(msg); user_msg_encode_err != nil {
			color.Set(color.FgRed)
			log.Println(user_msg_encode_err)
			color.Unset()
		}
	}
}
