package main

import (
	"log"
	"net"
	"os"
)

func main() {
	s := newServer()
	NewDBConn()
	NewRedisConn()
	NewRedisConnServer()
	defer sqliteHandler.Conn.Close()
	go s.run()

	PORT := os.Getenv("PORT")
	listener, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		log.Fatalf("unable to start server: %s", err.Error())
	}

	defer listener.Close()
	log.Printf("server started on :%s\n", PORT)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("failed to accept connection: %s", err.Error())
			continue
		}

		go s.newClient(conn)
	}
}
