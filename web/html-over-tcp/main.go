// Serve HTML over TCP.
package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	ln, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handle(conn)
	}
}

const page = `<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>Title</title>
	</head>
	<body>
		Body.
	</body>
</html>`

func handle(conn net.Conn) {
	defer conn.Close()
	if _, err := fmt.Fprint(conn, page); err != nil {
		log.Print(err)
	}
}
