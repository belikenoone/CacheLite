package main

import (
	"fmt"
	"net"
	"strings"
)

func main(){
	ln,err := net.Listen("tcp",":6380")
	if err != nil{
		panic(err)
	}
	fmt.Println("Server Started Accepting Connections")
	conn,err := ln.Accept()
	if err != nil{
		panic(err)
	}
	buffer:= make([]byte,1024)
	fmt.Println("Client Connected",conn.RemoteAddr().String())
	conn.Write([]byte("Connected Successfully\n"))

	store:= make(map[string]string)

	for {
		n,err := conn.Read(buffer)
		if err != nil{
			fmt.Println("Client Disconnected")
			return
		}
		clientMessage := strings.TrimSpace(string(buffer[:n]))
		fmt.Println("Client Send",clientMessage)

		parts:= strings.SplitN(clientMessage," ",3)
		command:=strings.ToUpper(parts[0])

		var key string
		var value string

		if len(parts) > 1{
			key = parts[1]
		}
		if len(parts) > 2{
			value = parts[2]
		}
		switch command{
			case "PING":
				conn.Write([]byte("PONG\n"))
			case "ECHO":
				echoParts:=strings.SplitN(clientMessage," ",2)
				if len(echoParts) < 2{
					conn.Write([]byte("What!!?\n"))
				}else{
					conn.Write([]byte(echoParts[1]+"\n"))
				}
			case "SET":
				if key == "" || value == ""{
					conn.Write([]byte("Invalid Arguments for set"))
					continue
				}
				store[key] = value
				conn.Write([]byte("OK!\n"))
			case "GET":
				val, exists := store[key]
				if !exists {
					conn.Write([]byte("(nil)\n"))
				} else {
					conn.Write([]byte(val + "\n"))
				}
			case "DEL":
				if key == ""{
					conn.Write([]byte("Provide a valid key name\n"))
					continue
				}
				_,exists := store[key]
				if !exists {
					conn.Write([]byte("0\n"))
				} else {
					delete(store,key)
					conn.Write([]byte("1\n"))
				}
			case "EXISTS":
				if key == ""{
					conn.Write([]byte("Provide a valid key name\n"))
					continue
				}
				_, exists := store[key]
				if !exists {
					conn.Write([]byte("0\n"))
				} else {
					conn.Write([]byte("1\n"))
				}
			case "KEYS":
				if len(store) == 0{
					conn.Write([]byte("empty\n"))
					continue
				}
				for k := range store{
					conn.Write([]byte(k+"\n"))
				}
			case "CLEAR":
				// If no keys, just say OK
				if len(store) == 0 {
					conn.Write([]byte("OK\n"))
					continue
				}

				for k := range store {	
					delete(store, k)
				}

    			conn.Write([]byte("OK\n"))

			case "QUIT":
				conn.Write([]byte("Connection Closed Successfully\n"))
				conn.Close()
				return

			default:
				conn.Write([]byte("INVALID ARGSS\n"))
		}
	}
	
}

