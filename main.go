package main

import (
	"log"
	"net"
)

func main() {
	s := newServer()
	newDBConn()
	defer sqliteHandler.Conn.Close()
	go s.run()

	listener, err := net.Listen("tcp4", "0.0.0.0:8888")
	if err != nil {
		log.Fatalf("unable to start server: %s", err.Error())
	}

	defer listener.Close()
	log.Printf("server started on :8888")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("failed to accept connection: %s", err.Error())
			continue
		}

		go s.newClient(conn)
	}
}
