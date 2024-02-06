package main

import (
	"client/lib"
)

func main() {
	client := lib.NewClient()

	client.ReadLoop()
	client.WriteLoop()

	defer func() {
		client.Conn.Close()
	}()
}
