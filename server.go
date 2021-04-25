package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

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

type Flashcard struct {
	fcID       int
	term       string
	definition string
	deckID     int
}

type Deck struct {
	deckID   int
	deckName string
	fcArray  *[]Flashcard
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
		case CMD_START:
			s.startGame(cmd.client)
		// case CMD_CREATEFC:
		// 	s.createfc(cmd.client, cmd.args[1], cmd.args[2])
		case CMD_CUSER:
			s.cuser(cmd.client)
		case CMD_SRD:
			s.setroomdeck(cmd.client, cmd.args[1])
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
		score:    0,
		vaild:    false,
		no_ques:  0,
	}

	c.readInput()
}

func newDBConn() {
	fmt.Println("Connecting to database...")
	db, err := sql.Open("mysql", "root:FgTQTzNM62cC63K@tcp(165.232.170.11:3306)/learninghub")
	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Println("Connected to database")
	sqliteHandler.Conn = db
}

func (s *server) nick(c *client, nick string) {

	c.msg(fmt.Sprint(nick))
	if len(nick) == 0 {
		c.msg(fmt.Sprintf("Invalid syntax"))
	}

	c.nick = nick
	c.msg(fmt.Sprintf("all right, I will call you %s", nick))
}

func (s *server) croom(c *client, roomName string) {

	c.msg(fmt.Sprintf("croom %s", roomName))
	r, ok := s.rooms[roomName]

	if !ok {

		ip := c.conn.RemoteAddr().String()
		r = &room{
			name:    roomName,
			members: make(map[net.Addr]*client),
			host:    ip, // ip of client
			status:  false,
			answer:  "",
			no_fc:   0,
		}

		s.rooms[roomName] = r

	} else {
		fmt.Println("Room already existed!...")
		return
	}

	r.members[c.conn.RemoteAddr()] = c

	s.quitCurrentRoom(c)
	c.room = r

	c.msg(fmt.Sprintf("welcome to %s", roomName))
	c.msg(fmt.Sprintf("if you want to play the game type ready"))
}

func (s *server) join(c *client, roomName string) {
	r, ok := s.rooms[roomName]
	if !ok {
		c.msg(fmt.Sprintf("Room is not available!"))
		return
	}
	if r.status == false {
		r.members[c.conn.RemoteAddr()] = c

		s.quitCurrentRoom(c)
		c.room = r

		r.broadcast(c, fmt.Sprintf("%s joined the room", c.nick))

		c.msg(fmt.Sprintf("welcome to %s", roomName))
		c.msg(fmt.Sprintf("if you want to play the game type ready"))
	} else {
		c.msg(fmt.Sprintf("game already start"))
	}

}

func (s *server) startGame(c *client) int {
	// BROADCAST QUESTION
	return 0
}

//////list all deck//////////////////////
func ListDecks(c *client) {
	sqlStatement := "select deckid, deckName from Deck_instance;"
	rows, err := sqliteHandler.Conn.Query(sqlStatement)
	if err != nil {
		return
	}
	for rows.Next() {
		var deckID string
		var deckName string
		err = rows.Scan(&deckID, &deckName)
		if err != nil {
			return
		}
		c.msg(fmt.Sprintf("%s : %s\n", deckID, deckName))
	}
	rows.Close()
}

////////////Check that deck name is already existed or not////////////////////////////
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
	rows.Close()
	return result
}

////////////time measure//////////////////////////
func timesup(c *client, msg string) {
	c.msg(msg)
	time.Sleep(3 * time.Second)
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
	rows.Close()
	return checkid, nil
}

func createDeck(c *client, deckname string) bool {
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

func createfc(c *client, listFC []Flashcard) {
	sqlStatement := `
		INSERT INTO Flashcard_instance(deckId,term,definition)
		VALUES(?,?,?)`
	for _, item := range listFC {
		_, err := sqliteHandler.Conn.Exec(sqlStatement, item.deckID, item.term, item.definition)
		if err != nil {
			return
		}
	}
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

func (s *server) cuser(c *client) {
	// see all number of user in the server or room
	return
}

func (s *server) setroomdeck(c *client, deckid string) {
	if c.conn.RemoteAddr().String() == c.room.host {
		deck_id, err := strconv.Atoi(deckid)
		if err != nil {
			c.msg(fmt.Sprintf("%s\n", err))
		}
		sqlStatement := `select Deck_instance.deckName, Deck_instance.deckId, 
		Flashcard_instance.flashcardID, Flashcard_instance.Term, 
		Flashcard_instance.definition from Flashcard_instance 
		inner join Deck_instance
		on Flashcard_instance.deckId = Deck_instance.deckId
		where Deck_instance.deckId = ?; `
		rows, err := sqliteHandler.Conn.Query(sqlStatement, deckid)
		if err != nil {
			return
		}
		var fcArray []Flashcard
		for rows.Next() {
			var tempFC Flashcard
			tempFC.deckID = deck_id
			err := rows.Scan(&c.room.deck.deckName, &c.room.deck.deckID,
				&tempFC.fcID, &tempFC.term, &tempFC.definition)
			fcArray = append(fcArray, tempFC)
			if err != nil {
				return
			}
		}
		rows.Close()
		c.room.no_fc = len(fcArray)
		c.room.deck.fcArray = &fcArray
		c.msg(fmt.Sprintf("This room have these shit\nDeckID=%d\nDeckName=%s\n", c.room.deck.deckID, c.room.deck.deckName))
		for _, item := range *c.room.deck.fcArray {
			c.msg(fmt.Sprintf("%d, %d, %s, %s\n", item.fcID, item.deckID, item.term, item.definition))
		}
		return
	} else {
		c.msg(fmt.Sprintf("You don't have the permission to change the room deck!"))
		return
	}
}
