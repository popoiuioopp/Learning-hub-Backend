package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"strings"
)

type server struct {
	rooms    map[string]*room
	commands chan command
}

type SQLHandler struct {
	Conn *sql.DB
}

var sqliteHandler SQLHandler

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
		}
	}
}

func (s *server) newClient(conn net.Conn) {
	log.Printf("new client has joined: %s", conn.RemoteAddr().String())

	c := &client{
		conn:     conn,
		nick:     "anonymous",
		status:   "0",
		commands: s.commands,
	}

	c.readInput()
}

func newDBConn() {

	db, err := sql.Open("mysql", "learninghub:FgTQTzNM62cC63K@tcp(139.59.106.148:3306)/learninghub")
	if err != nil {
		return
	}

	fmt.Println("Connected to database")
	sqliteHandler.Conn = db

	defer db.Close()
}

func (s *server) nick(c *client, nick string) {
	c.nick = nick
	c.msg(fmt.Sprintf("all right, I will call you %s", nick))
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
			status:    false,
		}
		s.rooms[roomName] = r
	} else {
		fmt.Println("Room already existed!...")
		return
	}

	r.members[c.conn.RemoteAddr()] = c

	s.quitCurrentRoom(c)
	c.room = r

	// r.broadcast(c, fmt.Sprintf("%s joined the room", c.nick))

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

	c.msg(fmt.Sprintf("welcome to %s", roomName))
}

func (s *server) startGame(c *client) {
	//start game
}

func checkDeckExist(db *sql.DB, name string) int {
	var result int
	statement := `SELECT COUNT(*) FROM learninghub.Deck_instance where learninghub.Deck_instance.deckName = ?;`
	rows, err := db.Query(statement, name)
	if err != nil {
		return 1
	}
	for rows.Next() {
		err = rows.Scan(&result)
		fmt.Print(result)
		if err != nil {
			return 1
		}
	}
	return result
}

func (s *server) createDeck(c *client, db *sql.DB, deckname string) {
	//create deck
	if checkDeckExist(db, deckname) == 0 {
		sqlStatement := `INSERT INTO Deck_instance(deckName, dateCreate) VALUES(?, NOW())`
		_, err := db.Exec(sqlStatement, deckname)

		if err != nil {
			return
		}

	} else {
		fmt.Println("Room already existed!...")
		return
	}
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
	}

	c.msg(fmt.Sprintf("available rooms: %s", strings.Join(rooms, ", ")))
}

func (s *server) msg(c *client, args []string) {
	msg := strings.Join(args[1:], " ")
	c.room.broadcast(c, c.nick+": "+msg)
}

func (s *server) quit(c *client) {
	log.Printf("client has left the chat: %s", c.conn.RemoteAddr().String())

	s.quitCurrentRoom(c)

	c.msg("sad to see you go =(")
	c.conn.Close()
}

func (s *server) quitCurrentRoom(c *client) {
	if c.room != nil {
		oldRoom := s.rooms[c.room.name]
		delete(s.rooms[c.room.name].members, c.conn.RemoteAddr())
		oldRoom.broadcast(c, fmt.Sprintf("%s has left the room", c.nick))
	}
}
