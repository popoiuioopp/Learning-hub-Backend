package main

import (
	"log"
	"net"
	"fmt"
)

func main() {
	s := newServer()
	NewDBConn()
	NewRedisConn()
	defer sqliteHandler.Conn.Close()
	go s.run()
	port := "33125"
	listener, err := net.Listen("tcp", fmt.Sprintf(":"+port))
	if err != nil {
		log.Fatalf("unable to start server: %s", err.Error())
	}

	defer listener.Close()
	log.Printf(fmt.Sprintf("server started on :"+ port))

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("failed to accept connection: %s", err.Error())
			continue
		}

		go s.newClient(conn)
	}
}
