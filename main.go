package main

import (
	"cloud-spanner/client"
	"cloud-spanner/server"
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		log.Panicln("Error in args")
	}

	arg := os.Args[1]

	if arg == "client" {
		fmt.Println("Launching client...")
		client.StartClient()
	} else {
		fmt.Println("Launching server...")
		server.StartServer()
	}
}
