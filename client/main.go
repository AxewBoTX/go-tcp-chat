package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/log"

	"client/lib"
)

func main() {
	// Getting The Username
	username_reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your username: ")
	username, username_read_err := username_reader.ReadString('\n')
	username = strings.TrimSpace(username)
	if username_read_err != nil {
		log.Fatal("Failed To Read Username", "Error", username_read_err)
	}
	if len(username) == 0 || username == "" {
		log.Fatal("Username is required!")
	}

	client := lib.NewClient()
	client.Username = username

	client.ReadLoop()
	client.WriteLoop()

	defer func() {
		client.Conn.Close()
	}()
}
