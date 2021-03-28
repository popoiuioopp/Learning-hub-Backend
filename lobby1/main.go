package main

import (
	"log"
	"net"
)

func main() {
	s := newServer()
	NewDBConn()
	NewRedisConn()
	defer sqliteHandler.Conn.Close()
	go s.run()

	listener, err := net.Listen("tcp", ":33125")
	if err != nil {
		log.Fatalf("unable to start server: %s", err.Error())
	}

	defer listener.Close()
	log.Printf("server started on :33125")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("failed to accept connection: %s", err.Error())
			continue
		}

		go s.newClient(conn)
	}
}
