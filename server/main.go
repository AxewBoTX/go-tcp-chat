package main

import (
	"server/lib"
)

func main() {
	server := lib.NewServer()
	server.Start()
}
