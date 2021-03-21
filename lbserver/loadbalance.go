package main

import (
	"flag"
	"io"
	"log"
	"net"
	"strings"
)

type Backends struct {
	servers []string
	n       int
}

func (b *Backends) Choose() string {
	idx := b.n % len(b.servers)
	b.n++
	return b.servers[idx]
}

func (b *Backends) String() string {
	return strings.Join(b.servers, ", ")
}

var (
	bind     = flag.String("bind", ":5000", "The address to bind on")
	balance  = flag.String("balance", "server1:5001,server2:5002", "The backend servers to balance connections across, separated by commas")
	backends *Backends
)

func init() {
	flag.Parse()

	if *bind == "" {
		log.Fatalln("specify the address to listen on with -bind")
	}

	servers := strings.Split(*balance, ",")
	if len(servers) == 1 && servers[0] == "" {
		log.Fatalln("please specify backend servers with -backends")
	}
	
	backends = &Backends{servers: servers}

}

func copy(wc io.WriteCloser, r io.Reader) {
	defer wc.Close()
	io.Copy(wc, r)
}

func handleConnection(userSide net.Conn, server string) {
	backendSide, err := net.Dial("tcp", server)
	if err != nil {
		userSide.Close()
		log.Printf("failed to dial %s: %s", server, err)
		return
	}

	go copy(backendSide, userSide)
	go copy(userSide, backendSide)

}

func main() {

	ln, err := net.Listen("tcp", *bind)

	if err != nil {
		log.Fatalf("failed to bind: %s", err)
	}

	log.Printf("listening on %s, balancing %s", *bind, backends)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("failed to accept: %s", err)
			continue
		}
		go handleConnection(conn, backends.Choose())
	}
}
