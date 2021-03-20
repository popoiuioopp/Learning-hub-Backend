package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

// global var
var totaluser int

type server struct {
	rooms    map[string]*room
	commands chan command
	// listener net.Listener
	// wg       sync.WaitGroup
	// quitt    chan interface{}
}

func newServer() *server {
	return &server{
		rooms:    make(map[string]*room),
		commands: make(chan command),
	}
}

func (s *server) run() {
	for cmd := range s.commands {
		switch cmd.id {
		case CMD_NICK:
			s.nick(cmd.client, cmd.args[1])
		case CMD_JOIN:
			s.join(cmd.client, cmd.args[1])
		case CMD_ROOMS:
			s.listRooms(cmd.client)
		case CMD_MSG:
			s.msg(cmd.client, cmd.args)
		case CMD_QUIT:
			s.quit(cmd.client)
		case CMD_CROOM:
			s.croom(cmd.client, cmd.args[1])
		case CMD_START:
			s.startGame(cmd.client)
		case CMD_CREATEFC:
			s.createfc(cmd.client, cmd.args[1], cmd.args[2])
		case CMD_LOGIN:
			s.login(cmd.client, cmd.args[1], cmd.args[2])
		case CMD_REGIS:
			s.regis(cmd.client, cmd.args[1], cmd.args[2], cmd.args[3])
		case CMD_TOTAL:
			s.total(cmd.client)
		}
	}
}

func (s *server) newClient(conn net.Conn) {

	// close(s.quitt)
	// s.listener.Close()
	// s.wg.Wait()
	log.Printf("new client has joined: %s", conn.RemoteAddr().String())
	totaluser += 1
	c := &client{
		conn:     conn,
		nick:     "anonymous",
		commands: s.commands,
	}
	c.msg(fmt.Sprintf("Welcome to Learning Hub!!!!!!!!"))
	c.readInput()
	// conn.Close()
}

func (s *server) nick(c *client, nick string) {
	c.nick = nick
	c.msg(fmt.Sprintf("all right, I will call you %s", nick))
	c.msg(fmt.Sprintf("Total user : %d", totaluser))
}

func (s *server) croom(c *client, roomName string) {

	c.msg(fmt.Sprintf("croom %s", roomName))
	r, ok := s.rooms[roomName]

	if !ok {

		c.msg(fmt.Sprintf("Select Your Flashcard..."))
		// Put Select Flashcard Function From Boss Here

		r = &room{
			name:      roomName,
			members:   make(map[net.Addr]*client),
			flashcard: "Flashcard", // Selected Flashcard
			host:      "host_id",   // Host ID
		}
		s.rooms[roomName] = r
	} else {
		fmt.Println("Room already existed!...")
		return
	}

	r.members[c.conn.RemoteAddr()] = c

	s.quitCurrentRoom(c)
	c.room = r

	r.broadcast(c, fmt.Sprintf("%s joined the room", c.nick))
	r.broadcast(c, fmt.Sprintf("Total player in current room : %d", len(r.members)))
	c.msg(fmt.Sprintf("welcome to %s", roomName))
}

func (s *server) join(c *client, roomName string) {
	r, ok := s.rooms[roomName]
	if !ok {

		c.msg(fmt.Sprintf("Room is not available!"))
		return
	}
	r.members[c.conn.RemoteAddr()] = c

	s.quitCurrentRoom(c)
	c.room = r

	r.broadcast(c, fmt.Sprintf("%s joined the room", c.nick))
	r.broadcast(c, fmt.Sprintf("Total player in the room : %d", len(r.members)))
	c.msg(fmt.Sprintf("welcome to %s", roomName))
}

func (s *server) startGame(c *client) {
	//start game
}

func (s *server) createfc(c *client, namefc string, total string) {
	//create fc
}

func (s *server) login(c *client, username string, pass string) {
	//login
}

func (s *server) regis(c *client, username string, pass string, verifypass string) {
	//regis
}

func (s *server) listRooms(c *client) {
	var rooms []string
	for name := range s.rooms {
		rooms = append(rooms, name)
		r := s.rooms[name]
		rooms = append(rooms, strconv.Itoa(len(r.members)))
		// c.msg(fmt.Sprintf("name : %d & %s", len(r.members), name))
	}
	c.msg(fmt.Sprintf("total rooms: %d", len(rooms)/2))
	c.msg(fmt.Sprintf("available rooms: %s", strings.Join(rooms, ", ")))

}

func (s *server) msg(c *client, args []string) {
	msg := strings.Join(args[1:len(args)], " ")
	c.room.broadcast(c, c.nick+": "+msg)
}

func (s *server) quit(c *client) {
	log.Printf("client has left the chat: %s", c.conn.RemoteAddr().String())

	s.quitCurrentRoom(c)

	c.msg("sad to see you go =(")
	c.conn.Close()
	// totaluser -= 1
}

func (s *server) quitCurrentRoom(c *client) {
	if c.room != nil {
		oldRoom := s.rooms[c.room.name]
		delete(s.rooms[c.room.name].members, c.conn.RemoteAddr())
		oldRoom.broadcast(c, fmt.Sprintf("%s has left the room", c.nick))
	}
}

func (s *server) total(c *client) {
	c.msg(fmt.Sprintf("Total user : %d", totaluser))
}
