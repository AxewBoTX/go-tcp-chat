package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"

	"client/lib"
)

func main() {
	// Getting The Username
	username_reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your username: ")
	username, username_read_err := username_reader.ReadString('\n')
	username = strings.TrimSpace(username)
	if username_read_err != nil {
		color.Set(color.FgRed)
		log.Fatal("Username Read Error:", username_read_err)
		color.Unset()
	}
	if len(username) == 0 || username == "" {
		color.Set(color.FgRed)
		log.Fatal("Username is required!")
		color.Unset()
	}

	client := lib.NewClient()
	client.Username = username

	client.ReadLoop()
	client.WriteLoop()

	defer func() {
		client.Conn.Close()
	}()
}
