package main

import (
	"fmt"
	"net"
	"strings"
)

func main() {
	// Start TCP server on port 6380
	ln, err := net.Listen("tcp", ":6380")
	if err != nil {
		panic(err) // Fatal: cannot start server
	}

	fmt.Println("Server Started. Accepting Connections...")

	// Accept only ONE client (for now)
	// Create the storage map for keys and values
	store := make(map[string]string)
	for{
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Cannot Accept Client Connection",err)
			continue
		}

		fmt.Println("Client Connected:", conn.RemoteAddr().String())


		// Handle this client's communication in a separate function
		go handleConnection(conn, store)

	}
}

func handleConnection(conn net.Conn, store map[string]string) {

	// Reusable byte buffer for reading incoming data
	buffer := make([]byte, 1024)
	// Send welcome message (optional)
	conn.Write([]byte("Connected Successfully\n"))

	for {
		// Read data from client (blocking call)
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Client Disconnected")
			return
		}

		// Convert bytes to string + trim whitespace/newline
		clientMessage := strings.TrimSpace(string(buffer[:n]))
		fmt.Println("Client Sent:", clientMessage)

		// Split into max 3 parts: command, key, value
		parts := strings.SplitN(clientMessage, " ", 3)
		command := strings.ToUpper(parts[0])

		var key string
		var value string

		// Extract key if present
		if len(parts) > 1 {
			key = parts[1]
		}

		// Extract value if present
		if len(parts) > 2 {
			value = parts[2]
		}

		// Command handling logic
		switch command {

		// Simple PING/PONG
		case "PING":
			conn.Write([]byte("PONG\n"))

		// ECHO command: returns the entire message after command
		case "ECHO":
			echoParts := strings.SplitN(clientMessage, " ", 2)
			if len(echoParts) < 2 {
				conn.Write([]byte("ERR missing message\n"))
			} else {
				conn.Write([]byte(echoParts[1] + "\n"))
			}

		// SET key value
		case "SET":
			if key == "" || value == "" {
				conn.Write([]byte("ERR syntax: SET key value\n"))
				continue
			}
			store[key] = value
			conn.Write([]byte("OK\n"))

		// GET key
		case "GET":
			val, exists := store[key]
			if !exists {
				conn.Write([]byte("(nil)\n"))
			} else {
				conn.Write([]byte(val + "\n"))
			}

		// DEL key
		case "DEL":
			if key == "" {
				conn.Write([]byte("ERR syntax: DEL key\n"))
				continue
			}
			_, exists := store[key]
			if !exists {
				conn.Write([]byte("0\n"))
			} else {
				delete(store, key)
				conn.Write([]byte("1\n"))
			}

		// EXISTS key
		case "EXISTS":
			if key == "" {
				conn.Write([]byte("ERR syntax: EXISTS key\n"))
				continue
			}
			_, exists := store[key]
			if exists {
				conn.Write([]byte("1\n"))
			} else {
				conn.Write([]byte("0\n"))
			}

		// KEYS — list all keys
		case "KEYS":
			if len(store) == 0 {
				conn.Write([]byte("(empty)\n"))
				continue
			}
			for k := range store {
				conn.Write([]byte(k + "\n"))
			}

		// CLEAR — delete all keys
		case "CLEAR":
			if len(store) == 0 {
				conn.Write([]byte("OK\n"))
				continue
			}
			// delete all keys
			for k := range store {
				delete(store, k)
			}
			conn.Write([]byte("OK\n"))

		// TYPE key — return type of stored value
		case "TYPE":
			if key == "" {
				conn.Write([]byte("ERR syntax: TYPE key\n"))
				continue
			}
			_, exists := store[key]
			if !exists {
				conn.Write([]byte("(nil)\n"))
			} else {
				conn.Write([]byte("string\n"))
			}

		// QUIT — close client connection
		case "QUIT":
			conn.Write([]byte("Connection Closed Successfully\n"))
			conn.Close()
			return

		// Any unknown command
		default:
			conn.Write([]byte("ERR unknown command\n"))
		}
	}
}
