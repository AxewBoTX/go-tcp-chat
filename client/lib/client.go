package lib

import (
	"bufio"
	"encoding/json"
	"net"
	"os"

	"github.com/charmbracelet/log"
)

func NewClient() *Client {
	conn, conn_err := net.Dial("tcp", SERVER_ADDR)
	if conn_err != nil {
		log.Fatal("Failed To Dial TCP Connection", "Error", conn_err)
	}
	return &Client{Conn: conn}
}

func (client *Client) ReadLoop() {
	// Client JOIN Message
	join_msg_encoder := json.NewEncoder(client.Conn)
	join_msg := Message{Client: *client, Method: "JOIN"}
	if join_msg_encoder_err := join_msg_encoder.Encode(join_msg); join_msg_encoder_err != nil {
		log.Fatal("Failed To Encode JOIN Message", "Error", join_msg_encoder_err)
	}
	go func() { // Running in a goroutine
		for {
			var server_msg Message
			server_msg_decoder := json.NewDecoder(client.Conn)
			if server_msg_decode_err := server_msg_decoder.Decode(&server_msg); server_msg_decode_err != nil {
				log.Fatal("Failed To Decode Server Message", "Error", server_msg_decode_err)
			}
			PrintMSG(server_msg)
		}
	}()
}

func (client *Client) WriteLoop() {
	user_input_scanner := bufio.NewScanner(os.Stdin)
	encoder := json.NewEncoder(client.Conn)
	for user_input_scanner.Scan() {
		msg := Message{Client: *client, Method: "MSG", Body: user_input_scanner.Text()}
		if user_msg_encode_err := encoder.Encode(msg); user_msg_encode_err != nil {
			log.Error("Failed To Encode User Message", "Error", user_msg_encode_err)
		}
	}
}
