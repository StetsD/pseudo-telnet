package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

var connections []net.Conn

func handleConnection(conn net.Conn) {
	defer func() {
		err := conn.Close()
		if err != nil {
			fmt.Printf("error during close connection: %v", err)
		}
	}()

	_, err := conn.Write([]byte(fmt.Sprintf("Welcome to chat mazafaka (%s)\n", conn.LocalAddr())))
	if err != nil {
		log.Fatalln(err)
	}

	connections = append(connections, conn)

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		msg := scanner.Text()
		log.Printf("Received from %s: %s", conn.RemoteAddr(), msg)
		if msg == "turboPower" {
			break
		}
		loc, _ := time.LoadLocation("UTC")
		now := time.Now().In(loc)
		for _, connection := range connections {
			if connection == conn {
				continue
			}

			go func(c net.Conn) {
				_, err := c.Write([]byte(fmt.Sprintf("%s : %s : %s\n", now, conn.RemoteAddr(), msg)))
				if err != nil {
					fmt.Printf("error happened during msg send %v\n", err)
				}
			}(connection)

		}
	}

	for i, connection := range connections {
		if connection == conn {
			connections = append(connections[:i], connections[i+1:]...)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error happend on connection with %s: %v", conn.RemoteAddr(), err)
	}

	fmt.Printf("Closing connection %s\n", conn.RemoteAddr())
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
