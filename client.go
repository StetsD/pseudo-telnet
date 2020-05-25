package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

type TelnetConnection struct {
}

func connectionReader(ctx context.Context, conn net.Conn, cancel context.CancelFunc) {
	data := bufio.NewScanner(conn)

outer:
	for {
		select {
		case <-ctx.Done():
			break outer
		default:
			if !data.Scan() {
				cancel()
				fmt.Println("can not read data")
				break outer
			}

			fmt.Printf("%s\n", data.Text())
		}
	}

	fmt.Println("Connection closed")
}

func connectionWriter(ctx context.Context, conn net.Conn, cancel context.CancelFunc) {
	msg := bufio.NewScanner(os.Stdin)

outer:
	for {
		select {
		case <-ctx.Done():
			break outer
		default:
			if msg.Scan() {
				_, err := conn.Write([]byte(msg.Text() + "\n"))

				if err != nil {
					fmt.Printf("error happened during msg send: %v", err)
					break outer
				}
			}
		}
	}

	fmt.Println("Connection closed")
}

func main() {
	dialer := net.Dialer{}
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Second)
	conn, err := dialer.DialContext(ctx, "tcp", "127.0.0.1:23")
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Fatalf("the error happened during connection close: %v", err)
		}
	}()

	if err != nil {
		log.Fatalln(err)
	}

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		connectionReader(ctx, conn, cancel)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		connectionWriter(ctx, conn, cancel)
		wg.Done()
	}()

	wg.Wait()
	conn.Close()
}
