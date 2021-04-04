package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"net"
	"fmt"
	//"strconv"
	"strings"
)

type Backends struct {
	servers []string
	n       int
}

func (b *Backends) Choose(idx int) string {
	// idx := b.n % len(b.servers)
	// b.n++
	// log.Printf("selected server idx: %d", idx-1)
	return b.servers[idx]
}

func (b *Backends) String() string {
	return strings.Join(b.servers, ", ")
}

func (b *Backends) HashRoom(name string) int {
	ascii := 0
	rune := []rune(name)
	for i := range rune {
		ascii += int(rune[i])
	}
	return ascii % len(b.servers)
}
   
var (
	bind = flag.String("bind", ":5000", "The address to bind on")
	balance  = flag.String("balance", "10.104.0.10:33125,10.104.0.3:33125,10.104.0.8:33125", "The backend servers to balance connections across, separated by commas")
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

func handleConnection(userSide net.Conn) {

	userSide.Write([]byte("> Please Select Your Room Name\n"))
	roomName, _ := bufio.NewReader(userSide).ReadString('\n')
	roomName = strings.TrimSuffix(roomName, "\r\n") // remove '\r\n'
	//log.Printf("room name: %s", roomName)
	idx := backends.HashRoom(roomName)
	//log.Printf("serverId: %d", idx)

	server := backends.Choose(idx)
	backendSide, err := net.Dial("tcp", server)

	if err != nil {
		userSide.Close()
		log.Printf("failed to dial %s: %s", server, err)
		return
	}

	// send roomname to server to create the room from the given name
	reader := bufio.NewReader(strings.NewReader(roomName))
	text, _ := reader.ReadString('\n')
	// send text to socket
	fmt.Fprint(backendSide, text+"\n")

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
		go handleConnection(conn)
	}
}
