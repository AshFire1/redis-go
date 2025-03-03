package main

import (
	"log"
	"net"
	"sync"
)

var cache sync.Map

func main() {
	listener, err := net.Listen("tcp", ":6380")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Listening on port 6380...")
	for {
		conn, err := listener.Accept()
		log.Print("New Connection")
		if err != nil {
			log.Println(err)
			continue
		}
		go startSession(conn)
	}
}

func startSession(conn net.Conn) {
	defer func() {
		log.Println("Clossing connection ", conn)
		conn.Close()
	}()
	defer func() {
		if err := recover(); err != nil {
			log.Println("Error in session:", err)
		}
	}()
	p := NewParser(conn)
	for {
		cmd, err := p.command()
		if err != nil {
			log.Println("Error parsing command:", err)
			conn.Write([]uint8("-ERR " + err.Error() + "\r\n"))
			return
		}
	}
}
