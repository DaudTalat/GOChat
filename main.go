package main

import (
	"fmt"
	"log"
	"net"
)

func main() {

	// init server
	server := init_server()
	go server.run()

	// open up listener
	ln, err := net.Listen("tcp", ":8080")
	if check_err(err) {
		fmt.Println("Error: Server Start Failed.")
	}

	log.Print("Started Server on PORT:8080")
	defer ln.Close()
	// keep checking for clients..
	for {
		cn, err := ln.Accept()

		if check_err(err) {
			fmt.Println("Error: Failed to Accept Connection...")
			continue
		}

		go server.init_client(cn) // init client - keep checking for messages
	}
}

func check_err(err error) bool {
	if err != nil {
		log.Fatalf("ERROR:%s", err)
		return true
	}

	return false
}
