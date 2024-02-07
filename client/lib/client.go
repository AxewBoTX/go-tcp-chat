package lib

import (
	"bufio"
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
	for user_input_scanner.Scan() {
		if _, send_err := client.Conn.Write([]byte(user_input_scanner.Text() + "\n")); send_err != nil {
			color.Set(color.FgRed)
			log.Println("User Input Read Error:", send_err)
			color.Unset()
			return
		}
	}
}
