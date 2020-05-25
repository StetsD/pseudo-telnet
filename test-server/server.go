package main

import (
	"fmt"
	"log"
	"net"
)

func handleConnection(conn net.Conn) {
	fmt.Println(conn.RemoteAddr())
}

func main() {
	gw, err := net.Listen("tcp", "127.0.0.1:23")

	defer func() {
		err := gw.Close()

		if err != nil {
			fmt.Printf("error happened during server stoping: %v", err)
		}
	}()

	if err != nil {
		log.Fatalf("server was interrupted with error: %v", err)
	}

	for {
		connection, err := gw.Accept()
		if err != nil {
			log.Fatalf("%v", err)
		}

		go handleConnection(connection)
	}
}
