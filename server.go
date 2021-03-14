package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type server struct {
	rooms    map[string]*room
	commands chan command
	num_user int 
}

type SQLHandler struct {
	Conn *sql.DB
}

type Deck struct {
	ID           int
	Name         string
	Capacity     int
	FlashCardsID []int
}

type Flashcard struct {
	ID         int
	term       string
	definition string
	deckID     int
}

type CreatingDeck struct {
	DeckStruct Deck
	Flashcards []Flashcard
}

var sqliteHandler SQLHandler

func newServer() *server {
	return &server{
		rooms:    make(map[string]*room),
		commands: make(chan command),
		num_user: 0,
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
		// case CMD_START:
		// 	s.startGame(cmd.client)
		case CMD_CREATEFC:
			s.createfc(cmd.client, cmd.args[1], cmd.args[2])
		case CMD_LOGIN:
			s.login(cmd.client, cmd.args[1], cmd.args[2])
		case CMD_REGIS:
			s.regis(cmd.client, cmd.args[1], cmd.args[2], cmd.args[3])
		case CMD_CUSER:
			s.cuser(cmd.client)
		}
	}
}

func (s *server) newClient(conn net.Conn) {
	log.Printf("new client has joined: %s", conn.RemoteAddr().String())
	s.num_user += 1

	c := &client{
		conn:     conn,
		nick:     "anonymous",
		status:   "0",
		commands: s.commands, 
	}
	
	c.readInput()
}

func newDBConn() {
	fmt.Println("Connecting to database...")
	db, err := sql.Open("mysql", "learninghub:FgTQTzNM62cC63K@tcp(143.198.204.127:3306)/learninghub")
	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Println("Connected to database")
	sqliteHandler.Conn = db
}

func (s *server) nick(c *client, nick string) {
	c.nick = nick
	c.msg(fmt.Sprintf("all right, I will call you %s", nick))
}

func (s *server) croom(c *client, roomName string) {

	c.msg(fmt.Sprintf("croom %s", roomName))
	r, ok := s.rooms[roomName]

	if !ok {

		//c.msg(fmt.Sprintf("Select Your Flashcard..."))
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
	c.msg(fmt.Sprintf("if you want to play the game type ready"))
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
	c.msg(fmt.Sprintf("if you want to play the game type ready"))
	
}

// func (s *server) startGame(c *client) int {
// 	// BROADCAST QUESTION 
// 	return
// }
//////list all deck//////////////////////
func ListDecks(c *client) {
	fmt.Print("Hello ListDeck")
	sqlStatement := "select deckid, deckName from Deck_instance;"
	rows, err := sqliteHandler.Conn.Query(sqlStatement)
	if err != nil {
		return
	}
	for rows.Next() {
		var deckID string
		var deckName string
		err = rows.Scan(&deckID, &deckName)
		c.msg(fmt.Sprintf("%s : %s\n", deckID, deckName))
	}
}
////////////check that deck name is already existed or not////////////////////////////
func checkDeckExist(deckname string) int {
	var result int
	statement := `SELECT COUNT(*) FROM learninghub.Deck_instance where learninghub.Deck_instance.deckName = ?;`
	rows, err := sqliteHandler.Conn.Query(statement, deckname)
	if err != nil {
		return 1
	}
	for rows.Next() {
		err = rows.Scan(&result)
		if err != nil {
			return 1
		}
	}
	return result
}

func checkDeckId(deckname string) (int, error) {
	var checkid int
	sqlStatement := `SELECT deckId FROM Deck_instance WHERE deckName = ? ORDER BY deckId DESC LIMIT 1 ` //check the lastest deckId and we will put it in the flashcard table
	rows, err := sqliteHandler.Conn.Query(sqlStatement, deckname)
	if err != nil {
		return 0, err
	}
	for rows.Next() {
		err = rows.Scan(&checkid)
		if err != nil {
			return 0, err
		}
	}
	return checkid, nil
}

func createDeck(c*client, deckname string) bool {
	if checkDeckExist(deckname) == 0 {
		sqlStatement := `INSERT INTO Deck_instance(deckName, dateCreate) VALUES(?, NOW())`
		_, err := sqliteHandler.Conn.Exec(sqlStatement, deckname)

		if err != nil {
			return false
		}
		c.msg(fmt.Sprintf("be able to create"))
		return true
	} else {
		c.msg(fmt.Sprintf("This Deck Name Already Used. You will be return to lobby"))
		return false
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
	c.msg(fmt.Sprintf("There are currently %d rooms", len(rooms)))
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
}

func (s *server) quitCurrentRoom(c *client) {
	if c.room != nil {
		oldRoom := s.rooms[c.room.name]
		delete(s.rooms[c.room.name].members, c.conn.RemoteAddr())
		oldRoom.broadcast(c, fmt.Sprintf("%s has left the room", c.nick))
	}
}

func (s *server) cuser(c *client){
	// try loop
	return
}