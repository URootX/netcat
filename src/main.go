package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	port = flag.String("p", "4444", "Port to listen on")
	host = flag.String("h", "0.0.0.0", "Host to bind to")
)

func main() {
	flag.Parse()

	// Create listener
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", *host, *port))
	if err != nil {
		log.Fatalf("Failed to listen on %s:%s: %v", *host, *port, err)
	}
	defer listener.Close()

	fmt.Printf("Listening on %s:%s...\n", *host, *port)

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\nShutting down...")
		listener.Close()
		os.Exit(0)
	}()

	// Accept connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		fmt.Printf("Connection received from %s\n", conn.RemoteAddr())

		// Handle each connection in a separate goroutine
		go handleConnection(conn)
	}
}

// handleConnection manages the bidirectional I/O mapping between stdin/stdout and the network connection
func handleConnection(conn net.Conn) {
	defer conn.Close()
	defer fmt.Printf("Connection from %s closed\n", conn.RemoteAddr())

	var wg sync.WaitGroup
	wg.Add(2)

	// Goroutine to copy from connection to stdout (remote -> local)
	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Panic in connection->stdout: %v", r)
			}
		}()

		_, err := io.Copy(os.Stdout, conn)
		if err != nil && err != io.EOF {
			log.Printf("Error copying from connection to stdout: %v", err)
		}
	}()

	// Goroutine to copy from stdin to connection (local -> remote)
	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Panic in stdin->connection: %v", r)
			}
		}()

		reader := bufio.NewReader(os.Stdin)
		for {
			// Read line from stdin
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Printf("Error reading from stdin: %v", err)
				break
			}

			// Write to connection
			_, err = conn.Write([]byte(line))
			if err != nil {
				log.Printf("Error writing to connection: %v", err)
				break
			}
		}
	}()

	// Wait for both goroutines to finish
	wg.Wait()
}
