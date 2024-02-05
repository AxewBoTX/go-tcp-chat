package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/fatih/color"

	"client/lib"
)

func main() {
	conn, conn_err := net.Dial("tcp", lib.SERVER_ADDR)
	if conn_err != nil {
		color.Set(color.FgRed)
		log.Fatal("TCP Connection Error:", conn_err)
		color.Unset()
	}

	go recieveData(conn)

	user_input_scanner := bufio.NewScanner(os.Stdin)
	for user_input_scanner.Scan() {
		if _, send_err := conn.Write([]byte(user_input_scanner.Text() + "\n")); send_err != nil {
			color.Set(color.FgRed)
			log.Println("User Input Read Error:", send_err)
			color.Unset()
			return
		}
	}

	defer func() {
		conn.Close()
	}()
}

func recieveData(conn net.Conn) {
	for {
		read_buff := make([]byte, 2048)
		if _, read_err := conn.Read(read_buff); read_err != nil {
			color.Set(color.FgRed)
			log.Println("Server Recieve Error:", read_err)
			color.Unset()
			return
		}
		fmt.Printf("-> %s\n", read_buff)
	}
}
