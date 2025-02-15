package main

import (
	"io"
	"log"
	"net"
)

var (
	addr         = "localhost:44304"
	upstreamAddr = "localhost:44300"
)

func main() {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	log.Printf("Listening on %s", addr)

	for {
		clientConn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		go handleConnection(clientConn)
	}
}

func handleConnection(clientConn net.Conn) {
	defer clientConn.Close()

	targetConn, err := net.Dial("tcp", upstreamAddr)
	if err != nil {
		log.Printf("Failed to connect to target: %v", err)
		return
	}
	defer targetConn.Close()

	go io.Copy(targetConn, clientConn)
	io.Copy(clientConn, targetConn)
}
